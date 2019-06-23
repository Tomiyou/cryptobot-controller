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
		build.Dir = botConfig.ArbitrageSrcPath
		output, err := build.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(output))

		// 2. run it with args
		test := exec.Command(
			botConfig.ArbitrageSrcPath+"crypto-arbitrage",
			"--config-path", "config/default_config.yaml",
			"--log-path", "csv",
			"--no-log",
		)
		test.Stdout = os.Stdout
		test.Stderr = os.Stderr
		if err = test.Run(); err != nil {
			fmt.Println(err)
			return
		}

		return
	},
}
