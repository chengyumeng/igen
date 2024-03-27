package cmd

import (
	"github.com/spf13/cobra"

	"github.com/chengyumeng/igen/pkg/decorate"
)

var (
	DecorateConfig = decorate.Config{}
	DecorateCmd    = &cobra.Command{
		Use:     "decorate",
		Short:   "Rely on the decorator pattern to encapsulate the function call standardized log",
		Long:    ``,
		Example: `igen decorate --help`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			d := decorate.Default(DecorateConfig)
			d.Decorate()
			return nil
		},
	}
)

func init() {
	DecorateCmd.Flags().StringVarP(&DecorateConfig.Source, "source", "s", "", "(source mode) Input Go source file,example: interface.go")
	DecorateCmd.Flags().StringVarP(&DecorateConfig.Destination, "destination", "d", "", "The file where the generated decorator code is stored")
	DecorateCmd.Flags().StringVarP(&DecorateConfig.ExplainFile, "explain", "e", "", "Type interpretation function mapping configuration")
	DecorateCmd.Flags().StringVarP(&DecorateConfig.SelfPackage, "self-package", "p", "", "The full package import path for the generated code")
	DecorateCmd.Flags().StringVarP(&DecorateConfig.Imports, "imports", "i", "", "Specify the imported package additionally")
	DecorateCmd.Flags().StringVarP(&DecorateConfig.Prom, "prom", "", "", "prometheus function template,example: prom.Create")
}
