package skills

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/skills/internal/config"
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

func TestDiscoverSubdir(t *testing.T) {
	dir := t.TempDir()

	// Create skills in <name>/SKILL.md format.
	mustMkdirAll(t, filepath.Join(dir, "go-rules"))
	mustWriteFile(t, filepath.Join(dir, "go-rules", "SKILL.md"), "---\nname: go-rules\ndescription: Go conventions\n---\n# go-rules\n\nGo conventions.\n")

	mustMkdirAll(t, filepath.Join(dir, "python"))
	mustWriteFile(t, filepath.Join(dir, "python", "SKILL.md"), "---\nname: python\ndescription: Python rules\n---\n# python\n\nPython rules.\n")

	// This directory has no SKILL.md, should be ignored.
	mustMkdirAll(t, filepath.Join(dir, "empty-dir"))

	// Non-md file, should be ignored.
	mustWriteFile(t, filepath.Join(dir, "not-a-skill.txt"), "ignored")

	cfg := &config.Config{SkillsDirs: []string{dir}}
	skills, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	names := map[string]bool{}
	for _, s := range skills {
		names[s.Name] = true
	}

	if !names["go-rules"] {
		t.Error("expected go-rules")
	}
	if !names["python"] {
		t.Error("expected python")
	}
}

func TestDiscoverFlatFiles(t *testing.T) {
	dir := t.TempDir()

	// Legacy flat file format.
	mustWriteFile(t, filepath.Join(dir, "go-rules.md"), "# go-rules\n\nGo conventions.\n")
	mustWriteFile(t, filepath.Join(dir, "python.md"), "# python\n\nPython rules.\n")
	mustWriteFile(t, filepath.Join(dir, "not-a-skill.txt"), "ignored")

	cfg := &config.Config{SkillsDirs: []string{dir}}
	skills, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	names := map[string]bool{}
	for _, s := range skills {
		names[s.Name] = true
	}

	if !names["go-rules"] {
		t.Error("expected go-rules")
	}
	if !names["python"] {
		t.Error("expected python")
	}
}

func TestDiscoverMixed(t *testing.T) {
	dir := t.TempDir()

	// Mix of subdir and flat file.
	mustMkdirAll(t, filepath.Join(dir, "dir-skill"))
	mustWriteFile(t, filepath.Join(dir, "dir-skill", "SKILL.md"), "---\nname: dir-skill\ndescription: From directory\n---\n# dir-skill\n")
	mustWriteFile(t, filepath.Join(dir, "flat-skill.md"), "# flat-skill\n\nFlat file.\n")

	cfg := &config.Config{SkillsDirs: []string{dir}}
	skills, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	names := map[string]bool{}
	for _, s := range skills {
		names[s.Name] = true
	}

	if !names["dir-skill"] {
		t.Error("expected dir-skill")
	}
	if !names["flat-skill"] {
		t.Error("expected flat-skill")
	}
}

func TestDiscoverPrecedence(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	mustMkdirAll(t, filepath.Join(dir1, "shared"))
	mustWriteFile(t, filepath.Join(dir1, "shared", "SKILL.md"), "---\nname: shared\ndescription: From dir1\n---\n# shared\n\nFrom dir1.\n")

	mustMkdirAll(t, filepath.Join(dir2, "shared"))
	mustWriteFile(t, filepath.Join(dir2, "shared", "SKILL.md"), "---\nname: shared\ndescription: From dir2\n---\n# shared\n\nFrom dir2.\n")

	mustMkdirAll(t, filepath.Join(dir2, "unique"))
	mustWriteFile(t, filepath.Join(dir2, "unique", "SKILL.md"), "---\nname: unique\ndescription: Only in dir2\n---\n# unique\n\nOnly in dir2.\n")

	cfg := &config.Config{SkillsDirs: []string{dir1, dir2}}
	skills, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	// The "shared" skill should come from dir1.
	for _, s := range skills {
		if s.Name == "shared" && s.Source != dir1 {
			t.Errorf("expected shared skill from %s, got %s", dir1, s.Source)
		}
	}
}

func TestDiscoverMissingDir(t *testing.T) {
	cfg := &config.Config{SkillsDirs: []string{"/nonexistent/path"}}
	skills, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 0 {
		t.Errorf("expected 0 skills, got %d", len(skills))
	}
}

func TestExtractFrontmatterField(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		field    string
		expected string
	}{
		{
			name:     "simple value",
			content:  "---\nname: go-rules\ndescription: Go conventions\n---\n# go-rules\n",
			field:    "description",
			expected: "Go conventions",
		},
		{
			name:     "quoted value",
			content:  "---\nname: go-rules\ndescription: \"Go conventions and best practices\"\n---\n",
			field:    "description",
			expected: "Go conventions and best practices",
		},
		{
			name:     "single quoted value",
			content:  "---\nname: go-rules\ndescription: 'Go conventions'\n---\n",
			field:    "description",
			expected: "Go conventions",
		},
		{
			name:     "name field",
			content:  "---\nname: my-skill\ndescription: test\n---\n",
			field:    "name",
			expected: "my-skill",
		},
		{
			name:     "no frontmatter",
			content:  "# Title\n\nNo frontmatter here.\n",
			field:    "description",
			expected: "",
		},
		{
			name:     "field not found",
			content:  "---\nname: test\n---\n# test\n",
			field:    "description",
			expected: "",
		},
		{
			name:     "empty content",
			content:  "",
			field:    "name",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFrontmatterField(tt.content, tt.field)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple",
			content:  "# Title\n\nThis is a description.\n",
			expected: "This is a description.",
		},
		{
			name:     "skip frontmatter",
			content:  "---\ntitle: foo\n---\n# Title\n\nDescription here.\n",
			expected: "Description here.",
		},
		{
			name:     "truncate long",
			content:  "# Title\n\n" + "This is a very long description line that exceeds eighty characters and should be truncated by the extraction function.",
			expected: "This is a very long description line that exceeds eighty characters and shoul...",
		},
		{
			name:     "empty",
			content:  "# Just a heading\n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDescription(tt.content)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestDiscoverFrontmatterDescription(t *testing.T) {
	dir := t.TempDir()

	mustMkdirAll(t, filepath.Join(dir, "with-frontmatter"))
	mustWriteFile(t, filepath.Join(dir, "with-frontmatter", "SKILL.md"), "---\nname: with-frontmatter\ndescription: Description from frontmatter\n---\n# with-frontmatter\n\nBody text here.\n")

	cfg := &config.Config{SkillsDirs: []string{dir}}
	skills, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}

	if skills[0].Description != "Description from frontmatter" {
		t.Errorf("expected description from frontmatter, got %q", skills[0].Description)
	}
}

func TestInstallAndCleanup(t *testing.T) {
	skillDir := t.TempDir()
	projectDir := t.TempDir()

	// Create a skill in directory format.
	mustMkdirAll(t, filepath.Join(skillDir, "test-skill"))
	skillPath := filepath.Join(skillDir, "test-skill", "SKILL.md")
	mustWriteFile(t, skillPath, "---\nname: test-skill\ndescription: Test skill\n---\n# test-skill\n\nTest skill.\n")

	// Create a .git marker so FindProjectRoot works.
	mustMkdirAll(t, filepath.Join(projectDir, ".git"))

	// Save and restore working directory.
	orig, _ := os.Getwd()
	defer mustChdir(t, orig)
	mustChdir(t, projectDir)

	cfg := &config.Config{
		SkillsDirs: []string{skillDir},
		Tools:      []string{"claude-code"},
	}

	available, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(available) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(available))
	}

	results, err := Install(cfg, available[0], false)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Err != nil {
		t.Fatalf("install failed: %v", r.Err)
	}
	if r.Skipped {
		t.Fatalf("install was skipped: %s", r.Reason)
	}

	// Verify the symlinked skill directory was created.
	expectedDest := filepath.Join(projectDir, ".claude", "skills", "test-skill")
	if r.Dest != expectedDest {
		t.Errorf("expected dest %q, got %q", expectedDest, r.Dest)
	}

	info, err := os.Lstat(r.Dest)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected symlink")
	}

	// Verify symlink target.
	target, err := os.Readlink(r.Dest)
	if err != nil {
		t.Fatal(err)
	}
	if target != filepath.Join(skillDir, "test-skill") {
		t.Errorf("expected target %q, got %q", filepath.Join(skillDir, "test-skill"), target)
	}

	// Verify Installed finds it.
	installed, err := Installed(cfg, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(installed) != 1 {
		t.Fatalf("expected 1 installed, got %d", len(installed))
	}
	if installed[0].Name != "test-skill" {
		t.Errorf("expected test-skill, got %s", installed[0].Name)
	}
	if !installed[0].IsSymlink {
		t.Error("expected IsSymlink to be true")
	}

	// Remove and verify.
	os.Remove(r.Dest)
	CleanupEmptyParent(r.Dest)

	// The skill entry should be cleaned up since it's gone.
	if _, err := os.Stat(filepath.Join(projectDir, ".claude", "skills", "test-skill")); err == nil {
		t.Error("expected skill subdir to be cleaned up")
	}
}

func TestInstallGlobalWithoutConfigFails(t *testing.T) {
	skillDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(skillDir, "test-skill"))
	mustWriteFile(t, filepath.Join(skillDir, "test-skill", "SKILL.md"), "---\nname: test-skill\ndescription: Test skill\n---\n# test-skill\n")

	cfg := &config.Config{
		SkillsDirs: []string{skillDir},
		Tools:      nil, // no tools configured
	}

	available, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Install(cfg, available[0], true)
	if err == nil {
		t.Fatal("expected error for global install without configured tools")
	}
}

func TestInstalledGlobalWithoutConfigFails(t *testing.T) {
	cfg := &config.Config{
		SkillsDirs: []string{t.TempDir()},
		Tools:      nil, // no tools configured
	}

	_, err := Installed(cfg, true)
	if err == nil {
		t.Fatal("expected error for global listing without configured tools")
	}
}

func TestInstallDefaultToolProject(t *testing.T) {
	skillDir := t.TempDir()
	projectDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(skillDir, "test-skill"))
	mustWriteFile(t, filepath.Join(skillDir, "test-skill", "SKILL.md"), "---\nname: test-skill\ndescription: Test skill\n---\n# test-skill\n")

	mustMkdirAll(t, filepath.Join(projectDir, ".git"))

	orig, _ := os.Getwd()
	defer mustChdir(t, orig)
	mustChdir(t, projectDir)

	cfg := &config.Config{
		SkillsDirs: []string{skillDir},
		Tools:      nil, // no tools configured, uses DefaultTool
	}

	available, err := Discover(cfg)
	if err != nil {
		t.Fatal(err)
	}

	results, err := Install(cfg, available[0], false)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Err != nil {
		t.Fatalf("install failed: %v", r.Err)
	}
	if r.Tool != "local" {
		t.Errorf("expected tool 'local', got %q", r.Tool)
	}

	expectedDest := filepath.Join(projectDir, ".agents", "skills", "test-skill")
	if r.Dest != expectedDest {
		t.Errorf("expected dest %q, got %q", expectedDest, r.Dest)
	}

	info, err := os.Lstat(r.Dest)
	if err != nil {
		t.Fatalf("expected symlink at %s: %v", r.Dest, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected symlink")
	}
}
