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
	if len(cfg.Tools) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", theme.Subtle.Render("all tools"))
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(cfg.Tools, ", "))
	}

	fmt.Fprintf(os.Stderr, "\n")

	// Supported tools.
	fmt.Fprintf(os.Stderr, "%s\n", theme.Bold.Render("Supported Tools:"))
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
		fmt.Fprintf(os.Stderr, "  %-12s %s  %s\n", t.Name, projectSupport, globalSupport)
	}

	return nil
}
