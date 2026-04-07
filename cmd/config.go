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
	Use:   "show",
	Short: "Show current configuration",
	RunE:  runConfigShow,
}

func init() {
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	fmt.Fprintf(os.Stderr, "%s\n", theme.Info.Render("Skills Configuration"))
	fmt.Fprintf(os.Stderr, "\n")

	// Skills directories.
	fmt.Fprintf(os.Stderr, "%s\n", theme.Bold.Render("Skills Directories:"))
	for _, dir := range cfg.SkillsDirs {
		exists := "  "
		if _, err := os.Stat(dir); err == nil {
			exists = theme.Success.Render("* ")
		} else {
			exists = theme.Warning.Render("! ")
		}
		fmt.Fprintf(os.Stderr, "  %s%s\n", exists, dir)
	}

	fmt.Fprintf(os.Stderr, "\n")

	// Tool filter.
	fmt.Fprintf(os.Stderr, "%s ", theme.Bold.Render("Tool Filter:"))
	if !tools.IsConfigured(cfg.Tools) {
		fmt.Fprintf(os.Stderr, "%s\n", theme.Subtle.Render("not set (using .agents/skills/ only)"))
		fmt.Fprintf(os.Stderr, "  %s\n", theme.Subtle.Render("Set SKILLS_TOOL to target specific agents, e.g.:"))
		fmt.Fprintf(os.Stderr, "  %s\n", theme.Subtle.Render("  export SKILLS_TOOL=\"claude-code,opencode\""))
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(cfg.Tools, ", "))
	}

	fmt.Fprintf(os.Stderr, "\n")

	// Active tools.
	fmt.Fprintf(os.Stderr, "%s\n", theme.Bold.Render("Active Tools:"))
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
		fmt.Fprintf(os.Stderr, "  %-18s %s  %s\n", t.Name, projectSupport, globalSupport)
	}

	// Show total available tools count when not all are configured.
	if !tools.IsConfigured(cfg.Tools) {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "%s %d agents available. Run %s to see all.\n",
			theme.Subtle.Render("hint:"),
			len(tools.All),
			theme.Bold.Render("SKILLS_TOOL=all skills config"))
	}

	return nil
}
