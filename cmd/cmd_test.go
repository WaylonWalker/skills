package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/skills/internal/skills"
	"github.com/spf13/cobra"
)

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path string, data string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustChdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
}

func TestRootCommandStructure(t *testing.T) {
	if rootCmd.Use != "skills" {
		t.Errorf("expected root use 'skills', got %q", rootCmd.Use)
	}

	// Verify subcommands are registered.
	expected := map[string]bool{
		"use": false, "list": false, "add": false, "remove": false, "config": false, "supported": false, "version": false,
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

func TestCommandAliases(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cobra.Command
		aliases []string
	}{
		{name: "add", cmd: addCmd, aliases: []string{"new", "create"}},
		{name: "use", cmd: useCmd, aliases: []string{"apply"}},
		{name: "supported", cmd: supportedCmd, aliases: []string{"tools", "agents"}},
		{name: "config show", cmd: configShowCmd, aliases: []string{"view"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, alias := range tt.aliases {
				if !containsAlias(tt.cmd.Aliases, alias) {
					t.Fatalf("expected alias %q on %s command", alias, tt.name)
				}
			}
		})
	}
}

func TestRootFlags(t *testing.T) {
	if rootCmd.PersistentFlags().Lookup("global") != nil {
		t.Fatal("did not expect root --global flag")
	}
	if rootCmd.PersistentFlags().Lookup("force") != nil {
		t.Fatal("did not expect root --force flag")
	}

	if gFlag := useCmd.Flags().Lookup("global"); gFlag == nil || gFlag.Shorthand != "g" {
		t.Fatal("expected --global flag on use command")
	}
	if gFlag := listCmd.Flags().Lookup("global"); gFlag == nil || gFlag.Shorthand != "g" {
		t.Fatal("expected --global flag on list command")
	}
	if gFlag := removeCmd.Flags().Lookup("global"); gFlag == nil || gFlag.Shorthand != "g" {
		t.Fatal("expected --global flag on remove command")
	}
	if fFlag := removeCmd.Flags().Lookup("force"); fFlag == nil || fFlag.Shorthand != "f" {
		t.Fatal("expected --force flag on remove command")
	}
	if addCmd.Flags().Lookup("global") != nil || addCmd.Flags().Lookup("force") != nil {
		t.Fatal("did not expect add command to expose global or force flags")
	}
	if configCmd.Flags().Lookup("global") != nil || configCmd.Flags().Lookup("force") != nil {
		t.Fatal("did not expect config command to expose global or force flags")
	}
	if supportedCmd.Flags().Lookup("global") != nil || supportedCmd.Flags().Lookup("force") != nil {
		t.Fatal("did not expect supported command to expose global or force flags")
	}
}

func TestListAliases(t *testing.T) {
	if !containsAlias(listCmd.Aliases, "ls") {
		t.Error("expected 'ls' alias on list command")
	}
}

func TestRemoveAliases(t *testing.T) {
	if !containsAlias(removeCmd.Aliases, "rm") {
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
	if err := executeArgs([]string{"add", "test-new-skill"}); err != nil {
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
	mustMkdirAll(t, filepath.Join(dir, "existing"))
	mustWriteFile(t, filepath.Join(dir, "existing", "SKILL.md"), "---\nname: existing\ndescription: test\n---\n# existing\n")

	err := executeArgs([]string{"add", "existing"})
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

func TestAddWithoutNameFailsNonInteractively(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "")

	err := runAdd(addCmd, nil)
	if err == nil {
		t.Fatal("expected non-interactive add to fail")
	}
	if !strings.Contains(err.Error(), "interactive input requires a terminal") {
		t.Fatalf("expected terminal guidance, got %q", err.Error())
	}
}

func TestUseUnknownSkillFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "claude-code")

	// Create one skill so discovery works (new dir format).
	mustMkdirAll(t, filepath.Join(dir, "real"))
	mustWriteFile(t, filepath.Join(dir, "real", "SKILL.md"), "---\nname: real\ndescription: A real skill\n---\n# real\n\nA real skill.\n")

	err := runUse(useCmd, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown skill")
	}
}

func TestUseInstallsSkill(t *testing.T) {
	skillDir := t.TempDir()
	projectDir := t.TempDir()

	// Create skill in new directory format.
	mustMkdirAll(t, filepath.Join(skillDir, "test"))
	mustWriteFile(t, filepath.Join(skillDir, "test", "SKILL.md"), "---\nname: test\ndescription: Test skill\n---\n# test\n\nTest skill.\n")
	mustMkdirAll(t, filepath.Join(projectDir, ".git"))

	t.Setenv("SKILLS_DIR", skillDir)
	t.Setenv("SKILLS_TOOL", "claude-code")

	orig, _ := os.Getwd()
	defer mustChdir(t, orig)
	mustChdir(t, projectDir)

	err := runUse(useCmd, []string{"test"})
	if err != nil {
		t.Fatalf("use command failed: %v", err)
	}

	// Verify the symlinked skill directory was created.
	link := filepath.Join(projectDir, ".claude", "skills", "test")
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
	t.Setenv("SKILLS_TOOL", "claude-code")

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

func TestRootDefaultsToUse(t *testing.T) {
	skillDir := t.TempDir()
	projectDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(skillDir, "test"))
	mustWriteFile(t, filepath.Join(skillDir, "test", "SKILL.md"), "---\nname: test\ndescription: Test skill\n---\n# test\n\nTest skill.\n")
	mustMkdirAll(t, filepath.Join(projectDir, ".git"))

	t.Setenv("SKILLS_DIR", skillDir)
	t.Setenv("SKILLS_TOOL", "claude-code")

	orig, _ := os.Getwd()
	defer mustChdir(t, orig)
	mustChdir(t, projectDir)

	if err := executeArgs([]string{"test"}); err != nil {
		t.Fatalf("root default use failed: %v", err)
	}

	link := filepath.Join(projectDir, ".claude", "skills", "test")
	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", link, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected symlink, got regular file")
	}
}

func TestVersionCommandWritesVersion(t *testing.T) {
	stdout := new(bytes.Buffer)
	versionCmd.SetOut(stdout)
	defer versionCmd.SetOut(nil)

	err := versionCmd.RunE(versionCmd, nil)
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}
	if strings.TrimSpace(stdout.String()) != version {
		t.Fatalf("expected version output %q, got %q", version, stdout.String())
	}
}

func TestNormalizeRootArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{name: "no args", args: nil, want: []string{"use"}},
		{name: "root subcommand", args: []string{"list"}, want: []string{"list"}},
		{name: "root alias", args: []string{"rm"}, want: []string{"rm"}},
		{name: "default skill name", args: []string{"go-rules"}, want: []string{"use", "go-rules"}},
		{name: "default global", args: []string{"-g", "go-rules"}, want: []string{"use", "-g", "go-rules"}},
		{name: "help", args: []string{"--help"}, want: []string{"--help"}},
		{name: "version flag", args: []string{"--version"}, want: []string{"--version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeRootArgs(tt.args)
			if strings.Join(got, "\x00") != strings.Join(tt.want, "\x00") {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
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

func TestRootHelpMentionsDefaultUseBehavior(t *testing.T) {
	if !strings.Contains(rootCmd.Long, "With no subcommand, skills behaves like 'skills use'.") {
		t.Fatal("expected root help to mention default use behavior")
	}
	if !strings.Contains(rootCmd.Example, "skills go-conventions") {
		t.Fatal("expected root examples to show implicit skill selection")
	}
}

func TestConfigShowRuns(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SKILLS_DIR", dir)
	t.Setenv("SKILLS_TOOL", "claude-code")

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	configShowCmd.SetOut(stdout)
	configShowCmd.SetErr(stderr)
	defer configShowCmd.SetOut(nil)
	defer configShowCmd.SetErr(nil)

	err := runConfigShow(configShowCmd, nil)
	if err != nil {
		t.Fatalf("config show failed: %v", err)
	}
	if stdout.Len() == 0 {
		t.Fatal("expected config output on stdout")
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no config output on stderr, got %q", stderr.String())
	}
}

func TestRemoveNonSymlinkFailsNonInteractivelyWithoutForce(t *testing.T) {
	skillDir := t.TempDir()
	projectDir := t.TempDir()
	t.Setenv("SKILLS_DIR", skillDir)
	t.Setenv("SKILLS_TOOL", "claude-code")

	mustMkdirAll(t, filepath.Join(projectDir, ".git"))
	installedDir := filepath.Join(projectDir, ".claude", "skills", "manual")
	if err := os.MkdirAll(installedDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(installedDir, "SKILL.md"), []byte("---\nname: manual\ndescription: test\n---\n# manual\n"), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	orig, _ := os.Getwd()
	defer mustChdir(t, orig)
	mustChdir(t, projectDir)

	if err := removeCmd.Flags().Set("global", "false"); err != nil {
		t.Fatalf("set global flag: %v", err)
	}
	if err := removeCmd.Flags().Set("force", "false"); err != nil {
		t.Fatalf("set force flag: %v", err)
	}
	defer func() {
		if err := removeCmd.Flags().Set("global", "false"); err != nil {
			t.Fatalf("reset global flag: %v", err)
		}
	}()
	defer func() {
		if err := removeCmd.Flags().Set("force", "false"); err != nil {
			t.Fatalf("reset force flag: %v", err)
		}
	}()

	err := runRemove(removeCmd, []string{"manual"})
	if err == nil {
		t.Fatal("expected non-interactive remove to fail")
	}
	if !strings.Contains(err.Error(), "rerun with --force") {
		t.Fatalf("expected force guidance, got %q", err.Error())
	}
	if _, statErr := os.Stat(installedDir); statErr != nil {
		t.Fatalf("expected skill to remain in place, got %v", statErr)
	}
}

func TestSupportedRuns(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe failed: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = origStdout
	}()

	err = runSupported(supportedCmd, nil)
	if err != nil {
		t.Fatalf("supported failed: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read failed: %v", err)
	}

	out := buf.String()
	if !containsString(out, "NAME") || !containsString(out, "PROJECT PATH") || !containsString(out, "GLOBAL PATH") {
		t.Fatalf("expected table headers in output, got %q", out)
	}
	if !containsString(out, "claude-code") {
		t.Fatalf("expected claude-code in output, got %q", out)
	}
	if !containsString(out, ".claude/skills") {
		t.Fatalf("expected project path in output, got %q", out)
	}
	if !containsString(out, ".config/opencode/skills") {
		t.Fatalf("expected opencode global path in output, got %q", out)
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

func containsAlias(aliases []string, want string) bool {
	for _, alias := range aliases {
		if alias == want {
			return true
		}
	}
	return false
}
