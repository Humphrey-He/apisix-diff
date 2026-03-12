package cli

import (
	"errors"
	"fmt"

	"github.com/awesomeProject/apidiff/internal/apisix"
	"github.com/awesomeProject/apidiff/internal/config"
	"github.com/awesomeProject/apidiff/internal/diff"
	"github.com/awesomeProject/apidiff/internal/render"
	"github.com/awesomeProject/apidiff/internal/validator"
	"github.com/spf13/cobra"
)

func newPlanCmd(opts *GlobalOptions) *cobra.Command {
	var filePath string
	var skipReachability bool
	var rulesPath string

	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Show diff between local config and live APISIX config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return errors.New("-f/--file is required")
			}

			ctx := cmd.Context()

			localCfg, err := config.LoadFile(filePath)
			if err != nil {
				return err
			}

			vOpts := validator.Options{SkipReachability: skipReachability, RulesPath: rulesPath}
			if err := validator.ValidateConfig(ctx, localCfg, vOpts); err != nil {
				return &ExitError{Code: 2, Err: err}
			}

			client := apisix.NewClient(opts.AdminURL, opts.Token, opts.Timeout)
			remoteCfg, err := client.FetchAll(ctx)
			if err != nil {
				return err
			}

			changes := diff.Compute(localCfg, remoteCfg)
			render.RenderPlan(cmd.OutOrStdout(), changes)

			if changes.HasChanges() {
				return &ExitError{Code: 1, Err: fmt.Errorf("diff detected")}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Local APISIX declarative config (YAML/JSON)")
	cmd.Flags().BoolVar(&skipReachability, "skip-reachability", false, "Skip upstream node reachability checks")
	cmd.Flags().StringVar(&rulesPath, "rules", "", "Rules file for plugin validation (YAML/JSON)")

	return cmd
}
