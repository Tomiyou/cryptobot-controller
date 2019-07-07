package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 1. build crypto-arbitrage
		build := exec.Command("go", "build")
		build.Dir = config.ArbitrageSrcPath
		err = build.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Built arbitrage bot.")

		// create csv folder if it doesn't exist
		_ = os.Mkdir("csv", 0777)

		// 2. run it with args
		test := exec.Command(
			config.ArbitrageSrcPath+"crypto-arbitrage",
			"--config-path", "config/default_config.yaml",
			"--log-path", "csv",
			"--no-log",
		)
		test.Env = append(test.Env, "CONTAINER_NAME=test_container")
		test.Stdout = os.Stdout
		test.Stderr = os.Stderr
		if err = test.Run(); err != nil {
			fmt.Println(err)
			return
		}

		return
	},
}
