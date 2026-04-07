package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/skills/internal/skills"
)

func TestRootCommandStructure(t *testing.T) {
	if rootCmd.Use != "skills" {
		t.Errorf("expected root use 'skills', got %q", rootCmd.Use)
	}

	// Verify subcommands are registered.
	expected := map[string]bool{
		"use": false, "list": false, "add": false, "remove": false, "config": false,
	}

	for _, sub := range rootCmd.Commands() {
		if _, ok := expected[sub.Name()]; ok {
			expected[sub.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("subcommand %q not registered on root", name)
		}
	}
}

func TestRootFlags(t *testing.T) {
	gFlag := rootCmd.PersistentFlags().Lookup("global")
	if gFlag == nil {
		t.Fatal("expected --global flag")
	}
	if gFlag.Shorthand != "g" {
		t.Errorf("expected shorthand 'g', got %q", gFlag.Shorthand)
	}

	fFlag := rootCmd.PersistentFlags().Lookup("force")
	if fFlag == nil {
		t.Fatal("expected --force flag")
	}
	if fFlag.Shorthand != "f" {
		t.Errorf("expected shorthand 'f', got %q", fFlag.Shorthand)
	}
}

func TestListAliases(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "list" {
			for _, alias := range sub.Aliases {
				if alias == "ls" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Error("expected 'ls' alias on list command")
	}
}

func TestRemoveAliases(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "remove" {
			for _, alias := range sub.Aliases {
				if alias == "rm" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Error("expected 'rm' alias on remove command")
	}
}

func TestSkillNames(t *testing.T) {
	tests := []struct {
		name     string
		input    []skills.Skill
		expected string
	}{
		{
			name:     "empty",
			input:    nil,
			expected: "",
		},
		{
			name:     "single",
			input:    []skills.Skill{{Name: "go-rules"}},
			expected: "go-rules",
		},
		{
			name: "multiple",
			input: []skills.Skill{
				{Name: "go-rules"},
				{Name: "python"},
				{Name: "typescript"},
			},
			expected: "go-rules, python, typescript",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skillNames(tt.input)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestAddCreatesSkillDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "")

	// Run the add command with a skill name.
	rootCmd.SetArgs([]string{"add", "test-new-skill"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	// Verify the directory and SKILL.md were created.
	path := filepath.Join(dir, "test-new-skill", "SKILL.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file at %s: %v", path, err)
	}

	// Check it contains frontmatter.
	got := string(content)
	if got == "" {
		t.Error("expected non-empty template content")
	}
	if !containsString(got, "name: test-new-skill") {
		t.Error("expected frontmatter to contain skill name")
	}
	if !containsString(got, "description:") {
		t.Error("expected frontmatter to contain description field")
	}
}

func TestAddDuplicateFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "")

	// Create existing skill directory.
	os.MkdirAll(filepath.Join(dir, "existing"), 0o755)
	os.WriteFile(filepath.Join(dir, "existing", "SKILL.md"), []byte("---\nname: existing\ndescription: test\n---\n# existing\n"), 0o644)

	rootCmd.SetArgs([]string{"add", "existing"})
	err := rootCmd.Execute()
	// SilenceErrors is true so Execute returns nil, but runAdd returns the error.
	// We need to call runAdd directly instead.
	_ = err

	err = runAdd(addCmd, []string{"existing"})
	if err == nil {
		t.Fatal("expected error for duplicate skill")
	}
}

func TestAddSanitizesName(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "")

	err := runAdd(addCmd, []string{"My Cool Skill.md"})
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}

	// Should be sanitized to my-cool-skill/SKILL.md
	path := filepath.Join(dir, "my-cool-skill", "SKILL.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected sanitized file at %s: %v", path, err)
	}
}

func TestAddEmptyNameFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "")

	err := runAdd(addCmd, []string{""})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUseUnknownSkillFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "claude")

	// Create one skill so discovery works (new dir format).
	os.MkdirAll(filepath.Join(dir, "real"), 0o755)
	os.WriteFile(filepath.Join(dir, "real", "SKILL.md"), []byte("---\nname: real\ndescription: A real skill\n---\n# real\n\nA real skill.\n"), 0o644)

	err := runUse(useCmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown skill")
	}
}

func TestUseInstallsSkill(t *testing.T) {
	skillDir := t.TempDir()
	projectDir := t.TempDir()

	// Create skill in new directory format.
	os.MkdirAll(filepath.Join(skillDir, "test"), 0o755)
	os.WriteFile(filepath.Join(skillDir, "test", "SKILL.md"), []byte("---\nname: test\ndescription: Test skill\n---\n# test\n\nTest skill.\n"), 0o644)
	os.MkdirAll(filepath.Join(projectDir, ".git"), 0o755)

	t.Setenv("SKILLS_DIR", skillDir)
	t.Setenv("SKILLS_TOOL", "claude")

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(projectDir)

	err := runUse(useCmd, []string{"test"})
	if err != nil {
		t.Fatalf("use command failed: %v", err)
	}

	// Verify the symlink was created.
	link := filepath.Join(projectDir, ".claude", "rules", "test.md")
	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", link, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected symlink, got regular file")
	}
}

func TestUseEmptySkillsDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "claude")

	// Should not error, just print a warning.
	err := runUse(useCmd, nil)
	if err != nil {
		t.Fatalf("expected nil error for empty dir, got: %v", err)
	}
}

func TestVersionIsSet(t *testing.T) {
	if rootCmd.Version == "" {
		t.Error("expected version to be set")
	}
}

func TestConfigSubcommand(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "config" {
			found = true
			// Check that show subcommand exists.
			hasShow := false
			for _, child := range sub.Commands() {
				if child.Name() == "show" {
					hasShow = true
				}
			}
			if !hasShow {
				t.Error("expected 'show' subcommand on config")
			}
		}
	}
	if !found {
		t.Error("expected 'config' subcommand on root")
	}
}

func TestConfigShowRuns(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "claude")

	err := runConfigShow(configShowCmd, nil)
	if err != nil {
		t.Fatalf("config show failed: %v", err)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
