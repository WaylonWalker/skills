// Package skills handles discovering, installing, and removing skill files.
package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/skills/internal/config"
	"github.com/WaylonWalker/skills/internal/tools"
)

// Skill represents a skill found in a skills directory.
type Skill struct {
	Name        string // skill name (directory name or filename without .md)
	Path        string // full path to the source SKILL.md or .md file
	DirPath     string // full path to the source skill directory when available
	Description string // from frontmatter or first non-header line
	Content     string // full file content
	Source      string // which skills directory it came from
}

// InstalledSkill represents a skill file currently installed for a tool.
type InstalledSkill struct {
	Name      string
	Path      string
	Tool      string
	IsSymlink bool
}

// InstallResult describes the outcome of installing a skill for one tool.
type InstallResult struct {
	Tool    string
	Dest    string
	Skipped bool
	Reason  string
	Err     error
}

// Discover finds all skills across all configured skills directories.
// It looks for both:
//   - <name>/SKILL.md directories (agentskills.io specification)
//   - <name>.md flat files (legacy format)
//
// Files in earlier directories take precedence over later ones with the same name.
func Discover(cfg *config.Config) ([]Skill, error) {
	seen := make(map[string]bool)
	var result []Skill

	for _, dir := range cfg.SkillsDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("reading %s: %w", dir, err)
		}

		for _, e := range entries {
			if e.IsDir() {
				// Check for <name>/SKILL.md pattern.
				name := e.Name()
				if seen[name] {
					continue
				}

				skillFile := filepath.Join(dir, name, "SKILL.md")
				content, err := os.ReadFile(skillFile)
				if err != nil {
					continue // no SKILL.md in this directory, skip
				}

				seen[name] = true
				desc := extractFrontmatterField(string(content), "description")
				if desc == "" {
					desc = extractDescription(string(content))
				}

				result = append(result, Skill{
					Name:        name,
					Path:        skillFile,
					DirPath:     filepath.Join(dir, name),
					Description: desc,
					Content:     string(content),
					Source:      dir,
				})
			} else if strings.HasSuffix(e.Name(), ".md") {
				// Legacy flat file format.
				name := strings.TrimSuffix(e.Name(), ".md")
				if seen[name] {
					continue
				}
				seen[name] = true

				path := filepath.Join(dir, e.Name())
				content, err := os.ReadFile(path)
				if err != nil {
					continue
				}

				desc := extractFrontmatterField(string(content), "description")
				if desc == "" {
					desc = extractDescription(string(content))
				}

				result = append(result, Skill{
					Name:        name,
					Path:        path,
					Description: desc,
					Content:     string(content),
					Source:      dir,
				})
			}
		}
	}

	return result, nil
}

// Install creates symlinks for a skill in the appropriate tool directories.
func Install(cfg *config.Config, skill Skill, global bool) ([]InstallResult, error) {
	if global && !tools.IsConfigured(cfg.Tools) {
		return nil, fmt.Errorf("global install requires explicit tool configuration\nSet SKILLS_TOOL to specify which tools to target, e.g.:\n  export SKILLS_TOOL=\"claude-code,opencode\"")
	}

	filtered := tools.Filtered(cfg.Tools)
	var results []InstallResult

	var projectRoot string
	if !global {
		var err error
		projectRoot, err = config.FindProjectRoot()
		if err != nil {
			return nil, fmt.Errorf("finding project root: %w", err)
		}
	}

	for _, tool := range filtered {
		var dest string
		if global {
			dest = tool.GlobalPath(skill.Name)
		} else {
			dest = tool.ProjectPath(projectRoot, skill.Name)
		}

		if dest == "" {
			scope := "project"
			if global {
				scope = "global"
			}
			results = append(results, InstallResult{
				Tool:    tool.Name,
				Skipped: true,
				Reason:  fmt.Sprintf("no %s support", scope),
			})
			continue
		}

		// Create parent directory.
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			results = append(results, InstallResult{
				Tool: tool.Name,
				Err:  fmt.Errorf("creating directory: %w", err),
			})
			continue
		}

		// If destination exists, handle it.
		if info, err := os.Lstat(dest); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				os.Remove(dest)
			} else {
				results = append(results, InstallResult{
					Tool:    tool.Name,
					Skipped: true,
					Reason:  "file exists and is not a symlink (use -f to overwrite)",
				})
				continue
			}
		}

		target := skill.Path
		if skill.DirPath != "" {
			target = skill.DirPath
		}

		if err := os.Symlink(target, dest); err != nil {
			results = append(results, InstallResult{
				Tool: tool.Name,
				Err:  err,
			})
			continue
		}

		results = append(results, InstallResult{
			Tool: tool.Name,
			Dest: dest,
		})
	}

	return results, nil
}

// Installed returns all skills currently installed for the given scope.
// All tools use the uniform <dir>/<name>/ pattern.
func Installed(cfg *config.Config, global bool) ([]InstalledSkill, error) {
	if global && !tools.IsConfigured(cfg.Tools) {
		return nil, fmt.Errorf("global listing requires explicit tool configuration\nSet SKILLS_TOOL to specify which tools to target, e.g.:\n  export SKILLS_TOOL=\"claude-code,opencode\"")
	}

	filtered := tools.Filtered(cfg.Tools)
	var result []InstalledSkill

	var projectRoot string
	if !global {
		var err error
		projectRoot, err = config.FindProjectRoot()
		if err != nil {
			return nil, fmt.Errorf("finding project root: %w", err)
		}
	}

	for _, tool := range filtered {
		var dir string
		if global {
			if tool.GlobalDir == "" {
				continue
			}
			dir = tool.GlobalDir
		} else {
			if tool.ProjectDir == "" {
				continue
			}
			dir = filepath.Join(projectRoot, tool.ProjectDir)
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, e := range entries {
			info, err := os.Lstat(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			if !e.IsDir() && info.Mode()&os.ModeSymlink == 0 {
				continue
			}
			skillPath := filepath.Join(dir, e.Name())
			resolvedInfo := info
			if info.Mode()&os.ModeSymlink != 0 {
				resolvedInfo, err = os.Stat(skillPath)
				if err != nil {
					continue
				}
			}
			if !resolvedInfo.IsDir() {
				continue
			}
			if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
				continue
			}
			result = append(result, InstalledSkill{
				Name:      e.Name(),
				Path:      skillPath,
				Tool:      tool.Name,
				IsSymlink: info.Mode()&os.ModeSymlink != 0,
			})
		}
	}

	return result, nil
}

// CleanupEmptyParent removes the parent directory if it is empty.
// This is useful after removing a per-skill directory entry.
func CleanupEmptyParent(path string) {
	dir := filepath.Dir(path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		os.Remove(dir)
	}
}

// extractFrontmatterField extracts a simple string value from YAML frontmatter.
// Handles frontmatter delimited by "---" lines. Returns empty string if not found.
func extractFrontmatterField(content, field string) string {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}

	for i := 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "---" {
			return "" // end of frontmatter, field not found
		}
		if strings.HasPrefix(trimmed, field+":") {
			val := strings.TrimPrefix(trimmed, field+":")
			val = strings.TrimSpace(val)
			// Strip surrounding quotes if present.
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}
			return val
		}
	}

	return ""
}

// extractDescription returns the first non-header, non-empty, non-frontmatter
// line from a markdown file, trimmed to 80 characters.
func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	inFrontmatter := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if len(trimmed) > 80 {
			return trimmed[:77] + "..."
		}
		return trimmed
	}

	return ""
}
