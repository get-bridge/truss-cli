package cmd

import (
	"fmt"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// BootstrapCmd represents the bootstrap command
var BootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap a new Truss deployment.",
	RunE: func(cmd *cobra.Command, args []string) error {
		set, err := cmd.Flags().GetStringToString("set")
		if err != nil {
			return err
		}
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		c, err := truss.LoadBootstrapConfig(config)
		if err != nil {
			return err
		}
		b, err := c.GetBootstrapper()
		defer b.Cleanup()
		if err != nil {
			return err
		}

		m := b.GetTemplateManifest()
		if m == nil {
			return errors.New("unable to load template manifest")
		}

		p := &truss.BootstrapParams{}
		p.LoadFromConfig(c)
		p.LoadFromFlags(set)

		if errs, err := p.Validate(m); err != nil {
			for _, err := range errs {
				fmt.Println(err)
			}
			return err
		}

		return b.Bootstrap(p)
	},
}

// BootstrapListTemplatesCmd represents the bootstrap list-templates command
var BootstrapListTemplatesCmd = &cobra.Command{
	Use:   "list-templates",
	Short: "List available templates for the bootstrap command.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, _ := cmd.Flags().GetString("config")
		c, err := truss.LoadBootstrapConfig(config)
		if err != nil {
			return err
		}
		ts, err := c.GetTemplateSource()
		if err != nil {
			return err
		}
		defer ts.Cleanup()

		t, err := ts.ListTemplates()
		if err != nil {
			return err
		}

		for _, t := range t {
			fmt.Println(t)
		}

		return nil
	},
}

func init() {
	BootstrapCmd.AddCommand(BootstrapListTemplatesCmd)
	BootstrapCmd.Flags().StringP("template", "t", "default", "Template to use")
	BootstrapCmd.Flags().StringToString("set", nil, "Set params on your template")
	BootstrapCmd.PersistentFlags().StringP("config", "f", "./bootstrap.truss.yaml", "Config file for bootstrapping")

	rootCmd.AddCommand(BootstrapCmd)
}
