package cmd

import (
	"github.com/spf13/cobra"
)

var (
	ShowVersion bool
	RootCmd     = &cobra.Command{
		Use:     "igen",
		Short:   "",
		Long:    ``,
		Example: `igen --help`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if ShowVersion {
				showVersion()
			}
			return nil
		},
	}
)

func init() {
	cobra.EnableCommandSorting = false
	RootCmd.Flags().BoolVarP(&ShowVersion, "version", "v", false, "show version")
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(DecorateCmd)
}
