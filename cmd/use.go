package cmd

import (
	"fmt"
	"os"

	"github.com/WaylonWalker/skills/internal/config"
	"github.com/WaylonWalker/skills/internal/skills"
	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/WaylonWalker/skills/internal/ui"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [skill-name]",
	Short: "Apply a skill to the current project or globally",
	Long: `Apply a skill by creating symlinks in tool-specific directories.

Without a skill name, opens a fuzzy picker to select from available skills.
Use -g to apply the skill globally instead of to the current project.`,
	Example: `  skills use                     Pick from available skills
  skills use go-conventions      Apply a specific skill
  skills use -g                  Pick and apply globally
  skills use -g go-conventions   Apply a specific skill globally`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUse,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	available, err := skills.Discover(cfg)
	if err != nil {
		return fmt.Errorf("discovering skills: %w", err)
	}

	if len(available) == 0 {
		fmt.Fprintln(os.Stderr, theme.Warning.Render("No skills found in: ")+cfg.SkillsDirDisplay())
		fmt.Fprintln(os.Stderr, theme.Subtle.Render("Run 'skills add' to create your first skill."))
		fmt.Fprintln(os.Stderr, theme.Subtle.Render("Set SKILLS_DIR to configure skill directories."))
		return nil
	}

	var selected *skills.Skill
	if len(args) == 1 {
		name := args[0]
		for i := range available {
			if available[i].Name == name {
				selected = &available[i]
				break
			}
		}
		if selected == nil {
			return fmt.Errorf("skill %q not found\nAvailable: %s", name, skillNames(available))
		}
	} else {
		picked, err := ui.Pick(available, "Select a skill to apply:")
		if err != nil {
			return err
		}
		if picked == nil {
			return nil
		}
		selected = picked
	}

	results, err := skills.Install(cfg, *selected, global)
	if err != nil {
		return err
	}

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: %s\n", theme.Error.Render("x"), r.Tool, r.Err)
		} else if r.Skipped {
			fmt.Fprintf(os.Stderr, "%s %s: %s\n", theme.Subtle.Render("-"), r.Tool, r.Reason)
		} else {
			fmt.Fprintf(os.Stderr, "%s %s -> %s\n", theme.Success.Render("*"), r.Tool, r.Dest)
		}
	}

	return nil
}

func skillNames(ss []skills.Skill) string {
	names := make([]string, len(ss))
	for i, s := range ss {
		names[i] = s.Name
	}
	result := ""
	for i, n := range names {
		if i > 0 {
			result += ", "
		}
		result += n
	}
	return result
}
