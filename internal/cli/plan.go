package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/awesomeProject/apidiff/internal/apisix"
	"github.com/awesomeProject/apidiff/internal/config"
	"github.com/awesomeProject/apidiff/internal/diff"
	"github.com/awesomeProject/apidiff/internal/render"
	"github.com/awesomeProject/apidiff/internal/validator"
	"github.com/spf13/cobra"
)

// newPlanCmd builds the "plan" command to show the diff between local and live config.
// It validates local config before calling the Admin API.
func newPlanCmd(opts *GlobalOptions) *cobra.Command {
	var filePath string
	var skipReachability bool
	var rulesPath string
	var colorMode string
	var output string

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

			out := cmd.OutOrStdout()
			switch strings.ToLower(output) {
			case "json":
				if err := render.RenderPlanJSON(out, changes); err != nil {
					return err
				}
			default:
				useColor, err := resolveColorMode(colorMode, out)
				if err != nil {
					return err
				}
				render.RenderPlan(out, changes, render.Options{Color: useColor})
			}

			if changes.HasChanges() {
				return &ExitError{Code: 1, Err: fmt.Errorf("diff detected")}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Local APISIX declarative config (YAML/JSON)")
	cmd.Flags().BoolVar(&skipReachability, "skip-reachability", false, "Skip upstream node reachability checks")
	cmd.Flags().StringVar(&rulesPath, "rules", "", "Rules file for plugin validation (YAML/JSON)")
	cmd.Flags().StringVar(&colorMode, "color", "auto", "Color output: auto, always, never")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format: text or json")

	return cmd
}

func resolveColorMode(mode string, out io.Writer) (bool, error) {
	switch strings.ToLower(mode) {
	case "auto":
		return isTerminal(out), nil
	case "always", "true", "1", "yes":
		return true, nil
	case "never", "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid --color value: %s", mode)
	}
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
