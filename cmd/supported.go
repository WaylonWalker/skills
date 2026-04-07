package cmd

import (
	"fmt"
	"os"
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
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tPROJECT PATH\tGLOBAL PATH")
	for _, tool := range tools.All {
		projectPath := tool.ProjectDir
		if projectPath == "" {
			projectPath = "-"
		}
		globalPath := tool.GlobalDir
		if globalPath == "" {
			globalPath = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", tool.Name, projectPath, globalPath)
	}
	return tw.Flush()
}
