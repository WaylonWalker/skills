package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage agent skills for AI coding assistants",
	Long: `skills manages your personal collection of agent skills.

Apply skills to projects or globally for tools like Claude, Copilot,
Cursor, OpenCode, Windsurf, Cline, Augment, and Roo Code.

Skills are directories containing a SKILL.md file with YAML frontmatter,
following the agentskills.io specification. Use this CLI to symlink them
into the directories each tool expects.

With no subcommand, skills behaves like 'skills use'.`,
	Example: `  skills                      Select and apply a skill to the current project
	  skills go-conventions       Apply a specific skill to the current project
	  skills use                  Select and apply a skill to the current project
	  skills use -g               Select and apply a skill globally
	  skills list                 Browse available skills with preview
	  skills add                  Create a new skill from a template
	  skills remove               Remove a skill from the current project
	  skills config               Show current configuration
	  skills version              Show the CLI version
	  skills supported            Show supported tools and paths`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() {
	if err := executeArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, theme.Error.Render("error: ")+err.Error())
		os.Exit(1)
	}
}

func executeArgs(args []string) error {
	rootCmd.SetArgs(normalizeRootArgs(args))
	return rootCmd.Execute()
}

func normalizeRootArgs(args []string) []string {
	if len(args) == 0 {
		return []string{"use"}
	}

	for _, arg := range args {
		if arg == "--help" || arg == "-h" || arg == "--version" || arg == "-v" {
			return args
		}
		if arg == "--" {
			break
		}
		if arg == "--global" || arg == "-g" {
			continue
		}
		if len(arg) > 0 && arg[0] == '-' {
			return args
		}
		if isRootSubcommand(arg) {
			return args
		}
		return append([]string{"use"}, args...)
	}

	return append([]string{"use"}, args...)
}

func isRootSubcommand(name string) bool {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == name || slices.Contains(cmd.Aliases, name) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.Version = version
}
