package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var DevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start cryptobot in developer mode.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		configFile, err := chooseConfigFile()
		if err != nil {
			return
		}

		// 1. build crypto-arbitrage
		build := exec.Command("go", "build")
		build.Dir = config.CryptobotSource
		err = build.Run()
		if err != nil {
			return
		}
		fmt.Println("Built arbitrage bot.")

		// create csv folder if it doesn't exist
		_ = os.Mkdir("csv", 0777)

		// 2. run it with args
		test := exec.Command(
			config.CryptobotSource+"cryptobot",
			"--config-path", "config/"+configFile,
			"--log-path", "csv",
			"--no-log",
		)
		test.Env = append(test.Env, "CONTAINER_NAME=test_container")
		test.Stdout = os.Stdout
		test.Stderr = os.Stderr
		if err = test.Run(); err != nil {
			return
		}

		return
	},
}
