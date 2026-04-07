package tools

import (
	"testing"
)

func TestFileName(t *testing.T) {
	tests := []struct {
		tool     Tool
		name     string
		expected string
	}{
		{
			tool:     Tool{Name: "claude"},
			name:     "my-skill",
			expected: "my-skill.md",
		},
		{
			tool: Tool{
				Name: "copilot",
				FileNameFunc: func(name string) string {
					return name + ".instructions.md"
				},
			},
			name:     "my-skill",
			expected: "my-skill.instructions.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tool.Name, func(t *testing.T) {
			got := tt.tool.FileName(tt.name)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestProjectPath(t *testing.T) {
	tool := Tool{
		Name:       "claude",
		ProjectDir: ".claude/rules",
	}

	got := tool.ProjectPath("/my/project", "go-rules")
	expected := "/my/project/.claude/rules/go-rules.md"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestProjectPathSubdir(t *testing.T) {
	tool := Tool{
		Name:        "opencode",
		NeedsSubdir: true,
		SubdirFile:  "SKILL.md",
		GlobalDir:   "/home/user/.config/opencode/skills",
	}

	// opencode has no project dir, so ProjectPath should return empty.
	got := tool.ProjectPath("/my/project", "go-rules")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}

	// But GlobalPath should work.
	got = tool.GlobalPath("go-rules")
	expected := "/home/user/.config/opencode/skills/go-rules/SKILL.md"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGlobalPathEmpty(t *testing.T) {
	tool := Tool{
		Name:       "cursor",
		ProjectDir: ".cursor/rules",
		// No GlobalDir set.
	}

	got := tool.GlobalPath("my-skill")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestFiltered(t *testing.T) {
	// Filter to just claude and copilot.
	result := Filtered([]string{"claude", "copilot"})

	if len(result) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(result))
	}

	names := make(map[string]bool)
	for _, r := range result {
		names[r.Name] = true
	}

	if !names["claude"] {
		t.Error("expected claude in results")
	}
	if !names["copilot"] {
		t.Error("expected copilot in results")
	}
}

func TestFilteredEmpty(t *testing.T) {
	// Empty filter returns all tools.
	result := Filtered(nil)
	if len(result) != len(All) {
		t.Errorf("expected %d tools, got %d", len(All), len(result))
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) != len(All) {
		t.Fatalf("expected %d names, got %d", len(All), len(names))
	}

	expected := map[string]bool{
		"claude": true, "copilot": true, "cursor": true, "opencode": true,
		"windsurf": true, "cline": true, "augment": true, "roo": true,
	}

	for _, name := range names {
		if !expected[name] {
			t.Errorf("unexpected tool name: %s", name)
		}
	}
}
