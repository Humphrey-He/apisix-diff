// Package cli defines the Cobra commands and exit-code behavior.
// It is responsible for wiring flags to validation, diff, and rendering.
package cli

import (
	"context"
	"time"

	"github.com/spf13/cobra"
)

// GlobalOptions holds flags shared across all subcommands.
// Values are passed into downstream clients and validators.
type GlobalOptions struct {
	// AdminURL is the APISIX Admin API base URL.
	AdminURL string
	// Token is the Admin API token used as X-API-KEY.
	Token string
	// Timeout is the HTTP timeout for Admin API calls.
	Timeout time.Duration
}

// NewRootCmd builds the root Cobra command and registers subcommands.
func NewRootCmd() *cobra.Command {
	opts := &GlobalOptions{}

	cmd := &cobra.Command{
		Use:           "apidiff",
		Short:         "APISIX declarative config diff and validation tool",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&opts.AdminURL, "admin-url", "http://127.0.0.1:9180", "APISIX Admin API base URL")
	cmd.PersistentFlags().StringVar(&opts.Token, "token", "", "APISIX Admin API token (X-API-KEY)")
	cmd.PersistentFlags().DurationVar(&opts.Timeout, "timeout", 5*time.Second, "HTTP timeout")

	cmd.AddCommand(newPlanCmd(opts))
	cmd.AddCommand(newValidateCmd(opts))
	cmd.AddCommand(newVersionCmd())

	return cmd
}

// Execute runs the CLI with a background context.
// Errors are handled by main to convert them to exit codes.
func Execute() error {
	return NewRootCmd().ExecuteContext(context.Background())
}
