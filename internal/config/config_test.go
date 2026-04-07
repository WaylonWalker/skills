package config

import (
	"os"
	"path/filepath"
	"testing"
)

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustChdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
}

func TestLoadDefaults(t *testing.T) {
	// Ensure env vars are cleared.
	t.Setenv("SKILLS_DIR", "")
	t.Setenv("SKILLS_TOOL", "")

	cfg := Load()

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "skills")

	if len(cfg.SkillsDirs) != 1 {
		t.Fatalf("expected 1 skills dir, got %d", len(cfg.SkillsDirs))
	}
	if cfg.SkillsDirs[0] != expected {
		t.Errorf("expected %q, got %q", expected, cfg.SkillsDirs[0])
	}
	if len(cfg.Tools) != 0 {
		t.Errorf("expected no tools filter, got %v", cfg.Tools)
	}
}

func TestLoadCustomDirs(t *testing.T) {
	t.Setenv("SKILLS_DIR", "/tmp/skills1, /tmp/skills2")
	t.Setenv("SKILLS_TOOL", "claude, copilot")

	cfg := Load()

	if len(cfg.SkillsDirs) != 2 {
		t.Fatalf("expected 2 skills dirs, got %d", len(cfg.SkillsDirs))
	}
	if cfg.SkillsDirs[0] != "/tmp/skills1" {
		t.Errorf("expected /tmp/skills1, got %q", cfg.SkillsDirs[0])
	}
	if cfg.SkillsDirs[1] != "/tmp/skills2" {
		t.Errorf("expected /tmp/skills2, got %q", cfg.SkillsDirs[1])
	}

	if len(cfg.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(cfg.Tools))
	}
	if cfg.Tools[0] != "claude" {
		t.Errorf("expected claude, got %q", cfg.Tools[0])
	}
	if cfg.Tools[1] != "copilot" {
		t.Errorf("expected copilot, got %q", cfg.Tools[1])
	}
}

func TestLoadExpandsHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	t.Setenv("SKILLS_DIR", "~/my-skills")
	t.Setenv("SKILLS_TOOL", "")

	cfg := Load()

	expected := filepath.Join(home, "my-skills")
	if len(cfg.SkillsDirs) != 1 || cfg.SkillsDirs[0] != expected {
		t.Errorf("expected %q, got %v", expected, cfg.SkillsDirs)
	}
}

func TestPrimarySkillsDir(t *testing.T) {
	cfg := &Config{SkillsDirs: []string{"/first", "/second"}}
	if got := cfg.PrimarySkillsDir(); got != "/first" {
		t.Errorf("expected /first, got %q", got)
	}
}

func TestFindProjectRoot(t *testing.T) {
	// Create a temp dir with a .git marker.
	dir := t.TempDir()
	sub := filepath.Join(dir, "a", "b")
	mustMkdirAll(t, sub)
	mustMkdirAll(t, filepath.Join(dir, ".git"))

	// Change to the subdirectory.
	orig, _ := os.Getwd()
	defer mustChdir(t, orig)
	mustChdir(t, sub)

	root, err := FindProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	if root != dir {
		t.Errorf("expected %q, got %q", dir, root)
	}
}

func TestLoadToolAll(t *testing.T) {
	t.Setenv("SKILLS_DIR", "")
	t.Setenv("SKILLS_TOOL", "all")

	cfg := Load()

	if len(cfg.Tools) != 1 || cfg.Tools[0] != "all" {
		t.Errorf("expected [all], got %v", cfg.Tools)
	}
}
