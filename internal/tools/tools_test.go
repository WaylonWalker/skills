package tools

import (
	"strings"
	"testing"
)

func TestProjectPath(t *testing.T) {
	tool := Tool{
		Name:       "claude-code",
		ProjectDir: ".claude/skills",
	}

	got := tool.ProjectPath("/my/project", "go-rules")
	expected := "/my/project/.claude/skills/go-rules/SKILL.md"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestProjectPathEmpty(t *testing.T) {
	tool := Tool{
		Name:      "custom",
		GlobalDir: "/home/user/.custom/skills",
		// No ProjectDir set.
	}

	got := tool.ProjectPath("/my/project", "go-rules")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestGlobalPath(t *testing.T) {
	tool := Tool{
		Name:       "opencode",
		ProjectDir: ".agents/skills",
		GlobalDir:  "/home/user/.config/opencode/skills",
	}

	got := tool.GlobalPath("go-rules")
	expected := "/home/user/.config/opencode/skills/go-rules/SKILL.md"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGlobalPathEmpty(t *testing.T) {
	tool := Tool{
		Name:       "no-global",
		ProjectDir: ".agents/skills",
		// No GlobalDir set.
	}

	got := tool.GlobalPath("my-skill")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestAllPathsAreSubdirs(t *testing.T) {
	// Every tool should produce paths ending in <name>/SKILL.md.
	for _, tool := range All {
		if tool.ProjectDir != "" {
			p := tool.ProjectPath("/root", "test-skill")
			if !strings.HasSuffix(p, "test-skill/SKILL.md") {
				t.Errorf("%s: ProjectPath should end with test-skill/SKILL.md, got %q", tool.Name, p)
			}
		}
		if tool.GlobalDir != "" {
			g := tool.GlobalPath("test-skill")
			if !strings.HasSuffix(g, "test-skill/SKILL.md") {
				t.Errorf("%s: GlobalPath should end with test-skill/SKILL.md, got %q", tool.Name, g)
			}
		}
	}
}

func TestFiltered(t *testing.T) {
	result := Filtered([]string{"claude-code", "opencode"})

	if len(result) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(result))
	}

	names := make(map[string]bool)
	for _, r := range result {
		names[r.Name] = true
	}

	if !names["claude-code"] {
		t.Error("expected claude-code in results")
	}
	if !names["opencode"] {
		t.Error("expected opencode in results")
	}
}

func TestFilteredEmpty(t *testing.T) {
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

	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	// Spot-check some known tools.
	for _, expected := range []string{"claude-code", "cursor", "opencode", "windsurf", "github-copilot"} {
		if !nameSet[expected] {
			t.Errorf("expected %q in tool names", expected)
		}
	}
}

func TestToolCount(t *testing.T) {
	if len(All) != 44 {
		t.Errorf("expected 44 tools, got %d", len(All))
	}
}
