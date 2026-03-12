package cli

import (
	"context"
	"time"

	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	AdminURL string
	Token    string
	Timeout  time.Duration
}

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

func Execute() error {
	return NewRootCmd().ExecuteContext(context.Background())
}
