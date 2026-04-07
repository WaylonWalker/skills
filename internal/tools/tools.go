// Package tools defines the supported AI coding assistant tools and their
// expected file paths for skill/rule/instruction files.
package tools

import (
	"os"
	"path/filepath"
	"runtime"
)

// Tool describes where a particular AI coding assistant expects its
// instruction files to be located.
type Tool struct {
	Name         string                        // tool identifier
	ProjectDir   string                        // relative to project root (empty = not supported)
	GlobalDir    string                        // absolute path (empty = not supported)
	FileNameFunc func(skillName string) string // custom filename; nil uses <name>.md
	NeedsSubdir  bool                          // if true, creates <dir>/<name>/<SubdirFile>
	SubdirFile   string                        // filename within subdirectory
}

// FileName returns the expected filename for a skill in this tool's directory.
func (t Tool) FileName(skillName string) string {
	if t.FileNameFunc != nil {
		return t.FileNameFunc(skillName)
	}
	return skillName + ".md"
}

// ProjectPath returns the full path where a skill should be installed at the
// project level. Returns empty string if the tool has no project-level support.
func (t Tool) ProjectPath(projectRoot, skillName string) string {
	if t.ProjectDir == "" {
		return ""
	}
	if t.NeedsSubdir {
		return filepath.Join(projectRoot, t.ProjectDir, skillName, t.SubdirFile)
	}
	return filepath.Join(projectRoot, t.ProjectDir, t.FileName(skillName))
}

// GlobalPath returns the full path where a skill should be installed globally.
// Returns empty string if the tool has no global support.
func (t Tool) GlobalPath(skillName string) string {
	if t.GlobalDir == "" {
		return ""
	}
	if t.NeedsSubdir {
		return filepath.Join(t.GlobalDir, skillName, t.SubdirFile)
	}
	return filepath.Join(t.GlobalDir, t.FileName(skillName))
}

// All is the registry of supported tools.
var All []Tool

func init() {
	home, _ := os.UserHomeDir()

	// Cline's global rules directory varies by platform.
	clineGlobal := filepath.Join(home, "Documents", "Cline", "Rules")
	if runtime.GOOS == "darwin" {
		clineGlobal = filepath.Join(home, "Documents", "Cline", "Rules")
	}

	All = []Tool{
		{
			Name:       "claude",
			ProjectDir: ".claude/rules",
			GlobalDir:  filepath.Join(home, ".claude", "rules"),
		},
		{
			Name:       "copilot",
			ProjectDir: ".github/instructions",
			FileNameFunc: func(name string) string {
				return name + ".instructions.md"
			},
		},
		{
			Name:       "cursor",
			ProjectDir: ".cursor/rules",
		},
		{
			Name:        "opencode",
			GlobalDir:   filepath.Join(home, ".config", "opencode", "skills"),
			NeedsSubdir: true,
			SubdirFile:  "SKILL.md",
		},
		{
			Name:       "windsurf",
			ProjectDir: ".windsurf/rules",
		},
		{
			Name:       "cline",
			ProjectDir: ".clinerules",
			GlobalDir:  clineGlobal,
		},
		{
			Name:       "augment",
			ProjectDir: ".augment/rules",
			GlobalDir:  filepath.Join(home, ".augment", "rules"),
		},
		{
			Name:       "roo",
			ProjectDir: ".roo/rules",
			GlobalDir:  filepath.Join(home, ".roo", "rules"),
		},
	}
}

// Filtered returns only the tools whose names appear in the given list.
// If names is empty, all tools are returned.
func Filtered(names []string) []Tool {
	if len(names) == 0 {
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

// Names returns the names of all supported tools.
func Names() []string {
	names := make([]string, len(All))
	for i, t := range All {
		names[i] = t.Name
	}
	return names
}
