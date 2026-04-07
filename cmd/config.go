package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/WaylonWalker/skills/internal/config"
	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/WaylonWalker/skills/internal/tools"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Long: `Display the current skills configuration including skills directories,
tool filter, and supported tools.`,
	Example: `  skills config        Show current configuration
  skills config show   Same as above`,
	RunE: runConfigShow,
}

var configShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"view"},
	Short:   "Show current configuration",
	RunE:    runConfigShow,
}

func init() {
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "%s\n", theme.Info.Render("Skills Configuration"))
	fmt.Fprintf(out, "\n")

	// Skills directories.
	fmt.Fprintf(out, "%s\n", theme.Bold.Render("Skills Directories:"))
	for _, dir := range cfg.SkillsDirs {
		exists := "  "
		if _, err := os.Stat(dir); err == nil {
			exists = theme.Success.Render("* ")
		} else {
			exists = theme.Warning.Render("! ")
		}
		fmt.Fprintf(out, "  %s%s\n", exists, dir)
	}

	fmt.Fprintf(out, "\n")

	// Tool filter.
	fmt.Fprintf(out, "%s ", theme.Bold.Render("Tool Filter:"))
	if !tools.IsConfigured(cfg.Tools) {
		fmt.Fprintf(out, "%s\n", theme.Subtle.Render("not set (using .agents/skills/ only)"))
		fmt.Fprintf(out, "  %s\n", theme.Subtle.Render("Set SKILLS_TOOL to target specific agents, e.g.:"))
		fmt.Fprintf(out, "  %s\n", theme.Subtle.Render("  export SKILLS_TOOL=\"claude-code,opencode\""))
	} else {
		fmt.Fprintf(out, "%s\n", strings.Join(cfg.Tools, ", "))
	}

	fmt.Fprintf(out, "\n")

	// Active tools.
	fmt.Fprintf(out, "%s\n", theme.Bold.Render("Active Tools:"))
	filtered := tools.Filtered(cfg.Tools)
	for _, t := range filtered {
		projectSupport := theme.Success.Render("project")
		if t.ProjectDir == "" {
			projectSupport = theme.Subtle.Render("--")
		}
		globalSupport := theme.Success.Render("global")
		if t.GlobalDir == "" {
			globalSupport = theme.Subtle.Render("--")
		}
		fmt.Fprintf(out, "  %-18s %s  %s\n", t.Name, projectSupport, globalSupport)
	}

	// Show total available tools count when not all are configured.
	if !tools.IsConfigured(cfg.Tools) {
		fmt.Fprintf(out, "\n")
		fmt.Fprintf(out, "%s %d agents available. Run %s to see all.\n",
			theme.Subtle.Render("hint:"),
			len(tools.All),
			theme.Bold.Render("SKILLS_TOOL=all skills config"))
	}

	return nil
}
