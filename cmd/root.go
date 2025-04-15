package cmd

import (
	"FinalProject/InformationScan/Modules"
	"FinalProject/cmd/flag_cli"
	"FinalProject/cmd/shell_cli"
	"FinalProject/common"
	"FinalProject/config"
	"fmt"
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
