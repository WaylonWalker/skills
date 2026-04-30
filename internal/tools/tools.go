// Package tools defines the supported AI coding assistant tools and their
// expected file paths for skill files.
//
// All tools follow the agentskills.io specification: skills are installed as
// <dir>/<name>/, with the directory containing SKILL.md and any optional
// companion files.
//
// The agent table is derived from https://github.com/vercel-labs/skills.
package tools

import (
	"os"
	"path/filepath"
)

// Tool describes where a particular AI coding assistant expects its
// skill files to be located. All tools use the uniform pattern
// <dir>/<name>/ for both project and global scopes.
type Tool struct {
	Name             string   // tool identifier (used with --agent flag / SKILLS_TOOL)
	ProjectDir       string   // primary project path, relative to project root (empty = not supported)
	GlobalDir        string   // primary global path, absolute (empty = not supported)
	ExtraProjectDirs []string // additional supported project paths
	ExtraGlobalDirs  []string // additional supported global paths
}

// ProjectPaths returns the full paths where a skill should be installed at the
// project level.
func (t Tool) ProjectPaths(projectRoot, skillName string) []string {
	var paths []string
	for _, dir := range append([]string{t.ProjectDir}, t.ExtraProjectDirs...) {
		if dir == "" {
			continue
		}
		paths = append(paths, filepath.Join(projectRoot, dir, skillName))
	}
	return uniqueStrings(paths)
}

// GlobalPaths returns the full paths where a skill should be installed
// globally.
func (t Tool) GlobalPaths(skillName string) []string {
	var paths []string
	for _, dir := range append([]string{t.GlobalDir}, t.ExtraGlobalDirs...) {
		if dir == "" {
			continue
		}
		paths = append(paths, filepath.Join(dir, skillName))
	}
	return uniqueStrings(paths)
}

// All is the registry of supported tools.
var All []Tool

func init() {
	home, _ := os.UserHomeDir()

	All = []Tool{
		// --- Agents with unique project paths ---
		{
			Name:       "claude-code",
			ProjectDir: ".claude/skills",
			GlobalDir:  filepath.Join(home, ".claude", "skills"),
		},
		{
			Name:       "windsurf",
			ProjectDir: ".windsurf/skills",
			GlobalDir:  filepath.Join(home, ".codeium", "windsurf", "skills"),
		},
		{
			Name:       "roo",
			ProjectDir: ".roo/skills",
			GlobalDir:  filepath.Join(home, ".roo", "skills"),
		},
		{
			Name:       "augment",
			ProjectDir: ".augment/skills",
			GlobalDir:  filepath.Join(home, ".augment", "skills"),
		},
		{
			Name:       "junie",
			ProjectDir: ".junie/skills",
			GlobalDir:  filepath.Join(home, ".junie", "skills"),
		},
		{
			Name:       "cody",
			ProjectDir: ".sourcegraph/skills",
			GlobalDir:  filepath.Join(home, ".sourcegraph", "skills"),
		},

		// --- Agents sharing .agents/skills project path ---
		{
			Name:       "cursor",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".cursor", "skills"),
		},
		{
			Name:       "github-copilot",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".copilot", "skills"),
			ExtraProjectDirs: []string{
				".github/skills",
			},
			ExtraGlobalDirs: []string{
				filepath.Join(home, ".agents", "skills"),
			},
		},
		{
			Name:       "opencode",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".config", "opencode", "skills"),
		},
		{
			Name:       "cline",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "codex",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".codex", "skills"),
		},
		{
			Name:       "aider",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".aider", "skills"),
		},
		{
			Name:       "void",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".void", "skills"),
		},
		{
			Name:       "pear",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".pear", "skills"),
		},
		{
			Name:       "zed",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".zed", "skills"),
		},
		{
			Name:       "continue",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".continue", "skills"),
		},
		{
			Name:       "goose",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".goose", "skills"),
		},
		{
			Name:       "trae",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".trae", "skills"),
		},
		{
			Name:       "aide",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".aide", "skills"),
		},
		{
			Name:       "qodo",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".qodo", "skills"),
		},
		{
			Name:       "tabnine",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".tabnine", "skills"),
		},
		{
			Name:       "gemini-cli",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".gemini", "skills"),
		},
		{
			Name:       "codeium",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".codeium", "skills"),
		},
		{
			Name:       "supermaven",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".supermaven", "skills"),
		},
		{
			Name:       "sourcegraph",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".sourcegraph-agent", "skills"),
		},

		// --- Agents sharing both .agents/skills project AND ~/.agents/skills global ---
		{
			Name:       "amp",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "kimi-cli",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "replit",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "universal",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "composio",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "devin",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "bolt",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "v0",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "lovable",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "stackblitz",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "same",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "softgen",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "cody-agent",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "idx",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "double",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "cloi",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "melty",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "manus",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
		{
			Name:       "hai",
			ProjectDir: ".agents/skills",
			GlobalDir:  filepath.Join(home, ".agents", "skills"),
		},
	}
}

// DefaultTool is the fallback when SKILLS_TOOL is not configured.
// It installs to the project-local .agents/skills/ directory only,
// with no global support. This is the safest default since .agents/skills/
// is the most widely shared project path across agents.
var DefaultTool = Tool{
	Name:       "local",
	ProjectDir: ".agents/skills",
	// No GlobalDir -- global installs require explicit tool configuration.
}

// Filtered returns only the tools whose names appear in the given list.
// If names is empty (no SKILLS_TOOL configured), only the DefaultTool is
// returned. This ensures the CLI does not scatter symlinks across dozens
// of tool directories without explicit opt-in.
// The special value "all" returns all tools.
func Filtered(names []string) []Tool {
	if len(names) == 0 {
		return []Tool{DefaultTool}
	}
	if len(names) == 1 && names[0] == "all" {
		return All
	}
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	var result []Tool
	for _, t := range All {
		if set[t.Name] {
			result = append(result, t)
		}
	}
	return result
}

// IsConfigured returns true if explicit tool names were provided
// (i.e. SKILLS_TOOL is set). When false, only the DefaultTool is active.
func IsConfigured(names []string) bool {
	return len(names) > 0
}

// Names returns the names of all supported tools.
func Names() []string {
	names := make([]string, len(All))
	for i, t := range All {
		names[i] = t.Name
	}
	return names
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}
