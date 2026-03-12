package cli

import (
	"errors"

	"github.com/awesomeProject/apidiff/internal/config"
	"github.com/awesomeProject/apidiff/internal/validator"
	"github.com/spf13/cobra"
)

// newValidateCmd builds the "validate" command for semantic checks only.
func newValidateCmd(opts *GlobalOptions) *cobra.Command {
	var filePath string
	var skipReachability bool
	var rulesPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate local APISIX declarative config",
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

			cmd.Println("validation passed")
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Local APISIX declarative config (YAML/JSON)")
	cmd.Flags().BoolVar(&skipReachability, "skip-reachability", false, "Skip upstream node reachability checks")
	cmd.Flags().StringVar(&rulesPath, "rules", "", "Rules file for plugin validation (YAML/JSON)")

	return cmd
}
