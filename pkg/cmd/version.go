package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/chengyumeng/igen/pkg/consts"
)

var (
	VersionCmd = &cobra.Command{
		Use:     "version",
		Short:   "",
		Long:    ``,
		Example: `igen version`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			showVersion()
			return nil
		},
	}
	versionData = `
 ___  ________  _______   ________       
|\  \|\   ____\|\  ___ \ |\   ___  \     
\ \  \ \  \___|\ \   __/|\ \  \\ \  \    
 \ \  \ \  \  __\ \  \_|/_\ \  \\ \  \   
  \ \  \ \  \|\  \ \  \_|\ \ \  \\ \  \   version   :  %s
   \ \__\ \_______\ \_______\ \__\\ \__\  language  :  %s
    \|__|\|_______|\|_______|\|__| \|__|  platform  :  %s/%s
`
)

func showVersion() {
	fmt.Printf(versionData, consts.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
