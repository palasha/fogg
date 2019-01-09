package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/errs"
	"github.com/chanzuckerberg/fogg/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	planCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	planCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
	planCmd.Flags().Bool("no-plugins", false, "do not apply fogg plugins; this may result in unexpected behavior.")
	rootCmd.AddCommand(planCmd)
}

var planCmd = &cobra.Command{
	Use:           "plan",
	Short:         "Run a plan",
	Long:          "plan will read fogg.json, use that to generate a plan and print that plan out. It will make no changes.",
	SilenceErrors: true, // If we don't silence here, cobra will print them. But we want to do that in cmd/root.go
	RunE: func(cmd *cobra.Command, args []string) error {
		setupDebug(debug)

		var e error
		// Set up fs
		pwd, e := os.Getwd()
		if e != nil {
			return errs.WrapUser(e, "can't get pwd")
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		verbose, e := cmd.Flags().GetBool("verbose")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse verbose flag")
		}

		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			return errs.WrapInternal(e, "couldn't parse config flag")
		}

		noPlugins, e := cmd.Flags().GetBool("no-plugins")
		if e != nil {
			return e
		}

		// check that we are at root of initialized git repo
		openGitOrExit(fs)

		config, err := readAndValidateConfig(fs, configFile, verbose)

		e = mergeConfigValidationErrors(err)
		if e != nil {
			return e
		}

		p, e := plan.Eval(config, verbose, noPlugins)
		if e != nil {
			return e
		}
		return plan.Print(p)
	},
}
