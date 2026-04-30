package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/WaylonWalker/skills/internal/tools"
	"github.com/spf13/cobra"
)

var supportedCmd = &cobra.Command{
	Use:     "supported",
	Aliases: []string{"tools", "agents"},
	Short:   "List supported tools and install paths",
	Long: `Print a table of supported tools and the paths where skills are installed
for project and global scopes.`,
	Example: `  skills supported`,
	RunE:    runSupported,
}

func init() {
	rootCmd.AddCommand(supportedCmd)
}

func runSupported(cmd *cobra.Command, args []string) error {
	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tPROJECT PATHS\tGLOBAL PATHS")
	for _, tool := range tools.All {
		projectPaths := joinPaths(append([]string{tool.ProjectDir}, tool.ExtraProjectDirs...))
		globalPaths := joinPaths(append([]string{tool.GlobalDir}, tool.ExtraGlobalDirs...))
		fmt.Fprintf(tw, "%s\t%s\t%s\n", tool.Name, projectPaths, globalPaths)
	}
	return tw.Flush()
}

func joinPaths(paths []string) string {
	var filtered []string
	for _, path := range paths {
		if path == "" {
			continue
		}
		filtered = append(filtered, path)
	}
	if len(filtered) == 0 {
		return "-"
	}
	return strings.Join(filtered, ", ")
}
