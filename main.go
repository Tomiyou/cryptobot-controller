package main

import (
	"fmt"

	"github.com/Tomiyou/cryptobot-controller/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crypto-arbitrage",
	Short: "Crypto-arbitrage bot written in go.",
}

func init() {
	// let the commands init first
	cmd.Initialize()

	// init cobra commands
	rootCmd.AddCommand(
		cmd.BuildCmd,
		cmd.StartCmd,
		cmd.StopCmd,
		cmd.UpdateCmd,
		cmd.EncryptCmd,
		cmd.DecryptCmd,
		cmd.CleanCmd,
		cmd.TestCmd,
	)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
