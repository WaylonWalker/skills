package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/skills/internal/config"
	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [skill-name]",
	Short: "Create a new skill from a template",
	Long: `Create a new skill directory with a SKILL.md file.

The skill follows the agentskills.io specification: a directory containing
a SKILL.md file with YAML frontmatter (name, description).

Without a name, prompts for one interactively. The skill is created
in the first directory listed in SKILLS_DIR.`,
	Example: `  skills add                  Prompted for a name
  skills add my-new-skill     Create with the given name`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	dir := cfg.PrimarySkillsDir()

	var name string
	if len(args) == 1 {
		name = args[0]
	} else {
		fmt.Fprint(os.Stderr, theme.Info.Render("Skill name: "))
		fmt.Scanln(&name)
		name = strings.TrimSpace(name)
	}

	if name == "" {
		return fmt.Errorf("skill name is required")
	}

	// Sanitize: lowercase, hyphens for spaces, strip extension.
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.TrimSuffix(name, ".md")

	skillDir := filepath.Join(dir, name)
	skillFile := filepath.Join(skillDir, "SKILL.md")

	if _, err := os.Stat(skillDir); err == nil {
		return fmt.Errorf("skill %q already exists at %s", name, skillDir)
	}

	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return fmt.Errorf("creating skill directory: %w", err)
	}

	template := fmt.Sprintf(`---
name: %s
description: Brief description of what this skill does
---

# %s

Add your instructions here.
`, name, name)

	if err := os.WriteFile(skillFile, []byte(template), 0o644); err != nil {
		// Clean up the directory on failure.
		os.Remove(skillDir)
		return fmt.Errorf("writing skill file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "%s Created %s\n", theme.Success.Render("*"), skillFile)

	if editor := os.Getenv("EDITOR"); editor != "" {
		fmt.Fprintf(os.Stderr, "%s Run '%s %s' to edit.\n", theme.Subtle.Render("->"), editor, skillFile)
	}

	return nil
}
