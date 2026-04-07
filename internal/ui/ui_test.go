package ui

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/skills/internal/skills"
)

func TestSkillItemInterface(t *testing.T) {
	s := skills.Skill{
		Name:        "go-rules",
		Description: "Go conventions and best practices",
		Path:        "/home/user/.config/skills/go-rules.md",
		Source:      "/home/user/git/public-skills",
		Content:     "# go-rules\n\nGo conventions.\n",
	}

	item := skillItem{skill: s}

	if got := item.Title(); got != "go-rules" {
		t.Errorf("Title() = %q, want %q", got, "go-rules")
	}

	if got := item.Description(); got != "Go conventions and best practices [public-skills]" {
		t.Errorf("Description() = %q, want %q", got, "Go conventions and best practices [public-skills]")
	}

	fv := item.FilterValue()
	if !strings.Contains(fv, "go-rules") {
		t.Errorf("FilterValue() should contain name, got %q", fv)
	}
	if !strings.Contains(fv, "Go conventions") {
		t.Errorf("FilterValue() should contain description, got %q", fv)
	}
	if !strings.Contains(fv, "public-skills") {
		t.Errorf("FilterValue() should contain source, got %q", fv)
	}
}

func TestSkillItemTitleMarksInstalled(t *testing.T) {
	item := skillItem{
		skill: skills.Skill{
			Name:        "go-rules",
			Description: "Go conventions and best practices",
			Source:      "/home/user/git/public-skills",
		},
		installed: true,
	}

	if got := item.Title(); got != "go-rules (installed)" {
		t.Errorf("Title() = %q, want %q", got, "go-rules (installed)")
	}
	if got := item.Description(); got != "Go conventions and best practices [public-skills]" {
		t.Errorf("Description() = %q, want %q", got, "Go conventions and best practices [public-skills]")
	}
}

func TestSkillItemDescriptionWithOnlySource(t *testing.T) {
	item := skillItem{skill: skills.Skill{Source: "/tmp/private-skills"}}
	if got := item.Description(); got != "private-skills" {
		t.Errorf("Description() = %q, want %q", got, "private-skills")
	}
}

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
	}{
		{
			name:    "simple heading",
			content: "# Hello World\n\nSome text.\n",
			width:   80,
		},
		{
			name:    "code block",
			content: "# Code\n\n```go\nfmt.Println(\"hello\")\n```\n",
			width:   80,
		},
		{
			name:    "narrow width clamps to 20",
			content: "# Narrow\n\nText.\n",
			width:   5,
		},
		{
			name:    "empty content",
			content: "",
			width:   80,
		},
		{
			name:    "bullet list",
			content: "# List\n\n- item one\n- item two\n- item three\n",
			width:   60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderMarkdown(tt.content, tt.width)
			// renderMarkdown should never return empty for non-empty input
			// (glamour always produces something).
			if tt.content != "" && strings.TrimSpace(result) == "" {
				t.Error("expected non-empty rendered output for non-empty input")
			}
		})
	}
}

func TestRenderMarkdownPreservesContent(t *testing.T) {
	content := "# My Skill\n\nThis is a test skill with specific content.\n"
	result := renderMarkdown(content, 80)

	// glamour wraps text with ANSI escape sequences, so we strip them
	// before checking for content.
	stripped := stripAnsi(result)
	if !strings.Contains(stripped, "My Skill") {
		t.Errorf("rendered output should contain heading text, got stripped: %q", stripped)
	}
	if !strings.Contains(stripped, "specific content") {
		t.Errorf("rendered output should contain body text, got stripped: %q", stripped)
	}
}

// stripAnsi removes ANSI escape sequences from a string for test assertions.
func stripAnsi(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// Skip until we find the terminator (a letter).
			j := i + 2
			for j < len(s) && !((s[j] >= 'A' && s[j] <= 'Z') || (s[j] >= 'a' && s[j] <= 'z')) {
				j++
			}
			if j < len(s) {
				j++ // skip the terminator
			}
			i = j
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}

func TestSkillItemEmptyFields(t *testing.T) {
	item := skillItem{skill: skills.Skill{}}

	if got := item.Title(); got != "" {
		t.Errorf("Title() for empty skill = %q, want empty", got)
	}
	if got := item.Description(); got != "" {
		t.Errorf("Description() for empty skill = %q, want empty", got)
	}
	if got := item.FilterValue(); got != " " {
		t.Errorf("FilterValue() for empty skill = %q, want ' '", got)
	}
}
