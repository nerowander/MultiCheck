package cmd

import (
	"fmt"
	"github.com/nerowander/MultiCheck/InformationScan/Modules"
	"github.com/nerowander/MultiCheck/cmd/flag_cli"
	"github.com/nerowander/MultiCheck/cmd/shell_cli"
	"github.com/nerowander/MultiCheck/common"
	"github.com/nerowander/MultiCheck/config"
	"os"
	"time"
)

func Execute() {
	startTime := time.Now()
	var info config.InfoScan
	args := os.Args[1:]
	if len(args) > 0 {
		flag_cli.Execute(&info)
		common.ParseInit(&info)
		Modules.HostScan(&info)
		common.GetSugestions()
		fmt.Printf("[*] Task finished, used time: %s\n", time.Since(startTime))
	} else {
		shell_cli.Start()
	}
}
