// Package config handles loading configuration from environment variables.
package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Config holds the resolved configuration for the CLI.
type Config struct {
	SkillsDirs []string // directories to search for skill files
	Tools      []string // tool filter (empty means project-local .agents/skills/ only)
}

// Load reads configuration from environment variables.
//
// SKILLS_DIR: comma-separated list of directories containing skill files.
// Defaults to ~/.config/skills.
//
// SKILLS_TOOL: comma-separated list of tools to target. Defaults to
// project-local .agents/skills/ only. Set to "all" for all tools.
func Load() *Config {
	cfg := &Config{}

	if env := os.Getenv("SKILLS_DIR"); env != "" {
		for _, d := range strings.Split(env, ",") {
			d = strings.TrimSpace(d)
			d = expandHome(d)
			if d != "" {
				cfg.SkillsDirs = append(cfg.SkillsDirs, d)
			}
		}
	} else {
		home, _ := os.UserHomeDir()
		cfg.SkillsDirs = []string{filepath.Join(home, ".config", "skills")}
	}

	if env := os.Getenv("SKILLS_TOOL"); env != "" {
		// "all" is a special value that explicitly targets all tools.
		if strings.TrimSpace(strings.ToLower(env)) == "all" {
			cfg.Tools = []string{"all"}
		} else {
			for _, t := range strings.Split(env, ",") {
				t = strings.TrimSpace(strings.ToLower(t))
				if t != "" {
					cfg.Tools = append(cfg.Tools, t)
				}
			}
		}
	}

	return cfg
}

// PrimarySkillsDir returns the first configured directory, used for creating
// new skills.
func (c *Config) PrimarySkillsDir() string {
	if len(c.SkillsDirs) > 0 {
		return c.SkillsDirs[0]
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "skills")
}

// SkillsDirDisplay returns a human-readable representation of the configured
// skills directories.
func (c *Config) SkillsDirDisplay() string {
	return strings.Join(c.SkillsDirs, ", ")
}

// FindProjectRoot walks up from the current working directory looking for
// common project root markers (.git, go.mod, package.json, etc.).
// Falls back to the current working directory if no marker is found.
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	markers := []string{
		".git",
		"go.mod",
		"package.json",
		"Cargo.toml",
		"pyproject.toml",
		"Makefile",
		"justfile",
		".project-root",
	}

	for {
		for _, m := range markers {
			if _, err := os.Stat(filepath.Join(dir, m)); err == nil {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return os.Getwd()
		}
		dir = parent
	}
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
