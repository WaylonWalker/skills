package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/WaylonWalker/skills/internal/config"
	"github.com/WaylonWalker/skills/internal/skills"
	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/WaylonWalker/skills/internal/tools"
	"github.com/WaylonWalker/skills/internal/ui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [skill-name]",
	Aliases: []string{"rm"},
	Short:   "Remove a skill from the current project or globally",
	Long: `Remove a skill by deleting symlinks from tool-specific directories.

If the installed directory is not a symlink, confirmation is required.
Use -f to force removal without confirmation.`,
	Example: `  skills remove                Remove a project skill (pick from installed)
  skills remove -g             Remove a global skill
  skills remove -f go-rules    Force remove even if not a symlink`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRemove,
}

func init() {
	removeCmd.Flags().BoolP("global", "g", false, "operate on global tool directories")
	removeCmd.Flags().BoolP("force", "f", false, "remove without confirmation")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	cfg := config.Load()

	installed, err := skills.Installed(cfg, global)
	if err != nil {
		return fmt.Errorf("listing installed skills: %w", err)
	}

	if len(installed) == 0 {
		scope := "project"
		if global {
			scope = "global"
		}
		fmt.Fprintf(os.Stderr, "%s No skills installed (%s).\n", theme.Warning.Render("!"), scope)
		if !tools.IsConfigured(cfg.Tools) {
			fmt.Fprintln(os.Stderr, theme.Subtle.Render("Only checking .agents/skills/ (default)."))
			fmt.Fprintln(os.Stderr, theme.Subtle.Render("Set SKILLS_TOOL to check tool-specific directories."))
		}
		return nil
	}

	var selected *skills.InstalledSkill
	if len(args) == 1 {
		name := args[0]
		for i := range installed {
			if installed[i].Name == name {
				selected = &installed[i]
				break
			}
		}
		if selected == nil {
			return fmt.Errorf("skill %q not found in installed skills", name)
		}
	} else {
		// Convert installed skills to Skill type for the picker.
		available := make([]skills.Skill, len(installed))
		for i, s := range installed {
			available[i] = skills.Skill{
				Name:        s.Name,
				Description: s.Tool + " -> " + s.Path,
				Path:        s.Path,
			}
		}

		picked, err := ui.Pick(available, "Select a skill to remove:", nil)
		if err != nil {
			return err
		}
		if picked == nil {
			return nil
		}

		for i := range installed {
			if installed[i].Name == picked.Name && installed[i].Path == picked.Path {
				selected = &installed[i]
				break
			}
		}
		if selected == nil {
			return nil
		}
	}

	// Warn if the installed directory is not a symlink and force is not set.
	if !selected.IsSymlink && !force {
		if !ui.IsInteractiveTerminal() {
			return fmt.Errorf("%s is not a symlink; rerun with --force to remove it non-interactively", selected.Path)
		}
		fmt.Fprintf(os.Stderr, "%s %s is not a symlink.\n", theme.Warning.Render("warning:"), selected.Path)
		fmt.Fprint(os.Stderr, "Remove anyway? [y/N] ")
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil && err.Error() != "unexpected newline" {
			return fmt.Errorf("reading confirmation: %w", err)
		}
		if !strings.HasPrefix(strings.ToLower(answer), "y") {
			fmt.Fprintln(os.Stderr, "Cancelled.")
			return nil
		}
	}

	if err := os.Remove(selected.Path); err != nil {
		return fmt.Errorf("removing %s: %w", selected.Path, err)
	}

	// Clean up empty parent directory (e.g. opencode subdirs).
	skills.CleanupEmptyParent(selected.Path)

	fmt.Fprintf(os.Stderr, "%s Removed %s (%s)\n", theme.Success.Render("*"), selected.Name, selected.Tool)
	return nil
}
