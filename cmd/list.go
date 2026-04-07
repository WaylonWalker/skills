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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Browse available skills with preview",
	Long: `List all available skills from your skills directories.

Opens an interactive browser with a markdown preview of each skill.
Use -g to show global installation status instead of project status.`,
	Example: `  skills list      Browse skills, see project install status
  skills list -g   Browse skills, see global install status`,
	Aliases: []string{"ls"},
	RunE:    runList,
}

func init() {
	listCmd.Flags().BoolP("global", "g", false, "show global installation status")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

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

	return ui.Preview(available, global, installedSkillNames(cfg, global))
}
