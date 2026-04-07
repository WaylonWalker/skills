package cmd

import (
	"fmt"
	"os"

	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	global  bool
	force   bool
)

var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage agent skills for AI coding assistants",
	Long: `skills manages your personal collection of agent skills.

Apply skills to projects or globally for tools like Claude, Copilot,
Cursor, OpenCode, Windsurf, Cline, Augment, and Roo Code.

Skills are directories containing a SKILL.md file with YAML frontmatter,
following the agentskills.io specification. Use this CLI to symlink them
into the directories each tool expects.`,
	Example: `  skills use                  Select and apply a skill to the current project
  skills use -g               Select and apply a skill globally
  skills list                 Browse available skills with preview
  skills add                  Create a new skill from a template
  skills remove               Remove a skill from the current project
  skills config               Show current configuration`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, theme.Error.Render("error: ")+err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "operate on global tool directories")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force the operation")
	rootCmd.Version = version
}
