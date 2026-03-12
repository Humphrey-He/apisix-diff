package cli

import (
	"github.com/spf13/cobra"
)

// version is overridden at build time via -ldflags when releasing.
var version = "dev"

// newVersionCmd prints the current build version.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version)
		},
	}
}
