package main

import (
	"fmt"
	"os"

	"github.com/chengyumeng/igen/pkg/cmd"
	"github.com/chengyumeng/igen/pkg/consts"
)

// igen version
var (
	Version = "0.6.0"
)

func main() {
	consts.Version = Version
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
