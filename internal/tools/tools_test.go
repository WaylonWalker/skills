package tools

import (
	"strings"
	"testing"
)

func TestProjectPaths(t *testing.T) {
	tool := Tool{
		Name:       "claude-code",
		ProjectDir: ".claude/skills",
	}

	got := tool.ProjectPaths("/my/project", "go-rules")
	if len(got) != 1 {
		t.Fatalf("expected 1 path, got %d", len(got))
	}
	expected := "/my/project/.claude/skills/go-rules"
	if got[0] != expected {
		t.Errorf("expected %q, got %q", expected, got[0])
	}
}

func TestProjectPathsEmpty(t *testing.T) {
	tool := Tool{
		Name:      "custom",
		GlobalDir: "/home/user/.custom/skills",
		// No ProjectDir set.
	}

	got := tool.ProjectPaths("/my/project", "go-rules")
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestGlobalPaths(t *testing.T) {
	tool := Tool{
		Name:       "opencode",
		ProjectDir: ".agents/skills",
		GlobalDir:  "/home/user/.config/opencode/skills",
	}

	got := tool.GlobalPaths("go-rules")
	if len(got) != 1 {
		t.Fatalf("expected 1 path, got %d", len(got))
	}
	expected := "/home/user/.config/opencode/skills/go-rules"
	if got[0] != expected {
		t.Errorf("expected %q, got %q", expected, got[0])
	}
}

func TestGlobalPathsEmpty(t *testing.T) {
	tool := Tool{
		Name:       "no-global",
		ProjectDir: ".agents/skills",
		// No GlobalDir set.
	}

	got := tool.GlobalPaths("my-skill")
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestCopilotPathsIncludeAllDocumentedLocations(t *testing.T) {
	var copilot Tool
	for _, tool := range All {
		if tool.Name == "github-copilot" {
			copilot = tool
			break
		}
	}

	projectPaths := copilot.ProjectPaths("/repo", "test-skill")
	if len(projectPaths) != 2 {
		t.Fatalf("expected 2 project paths, got %d", len(projectPaths))
	}
	if projectPaths[0] != "/repo/.agents/skills/test-skill" {
		t.Errorf("unexpected first project path: %q", projectPaths[0])
	}
	if projectPaths[1] != "/repo/.github/skills/test-skill" {
		t.Errorf("unexpected second project path: %q", projectPaths[1])
	}

	globalPaths := copilot.GlobalPaths("test-skill")
	if len(globalPaths) != 2 {
		t.Fatalf("expected 2 global paths, got %d", len(globalPaths))
	}
	if !strings.HasSuffix(globalPaths[0], "/.copilot/skills/test-skill") {
		t.Errorf("unexpected first global path: %q", globalPaths[0])
	}
	if !strings.HasSuffix(globalPaths[1], "/.agents/skills/test-skill") {
		t.Errorf("unexpected second global path: %q", globalPaths[1])
	}
}

func TestAllPathsAreSubdirs(t *testing.T) {
	// Every tool should produce paths ending in <name>/.
	for _, tool := range All {
		for _, p := range tool.ProjectPaths("/root", "test-skill") {
			if !strings.HasSuffix(p, "test-skill") {
				t.Errorf("%s: ProjectPaths should end with test-skill, got %q", tool.Name, p)
			}
		}
		for _, g := range tool.GlobalPaths("test-skill") {
			if !strings.HasSuffix(g, "test-skill") {
				t.Errorf("%s: GlobalPaths should end with test-skill, got %q", tool.Name, g)
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
	// When no tools are configured, only the DefaultTool is returned.
	result := Filtered(nil)
	if len(result) != 1 {
		t.Errorf("expected 1 tool (DefaultTool), got %d", len(result))
	}
	if result[0].Name != "local" {
		t.Errorf("expected default tool name 'local', got %q", result[0].Name)
	}
	if result[0].ProjectDir != ".agents/skills" {
		t.Errorf("expected ProjectDir '.agents/skills', got %q", result[0].ProjectDir)
	}
	if result[0].GlobalDir != "" {
		t.Errorf("expected empty GlobalDir, got %q", result[0].GlobalDir)
	}
}

func TestFilteredAll(t *testing.T) {
	// The special value "all" returns every tool.
	result := Filtered([]string{"all"})
	if len(result) != len(All) {
		t.Errorf("expected %d tools, got %d", len(All), len(result))
	}
}

func TestIsConfigured(t *testing.T) {
	if IsConfigured(nil) {
		t.Error("expected false for nil")
	}
	if IsConfigured([]string{}) {
		t.Error("expected false for empty slice")
	}
	if !IsConfigured([]string{"claude-code"}) {
		t.Error("expected true for non-empty slice")
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
