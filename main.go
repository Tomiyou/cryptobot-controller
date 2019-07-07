package main

import (
	"fmt"
	"log"

	"github.com/Tomiyou/cryptobot-controller/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crypto-arbitrage",
	Short: "Crypto-arbitrage bot written in go.",
}

func main() {
	rootCmd.AddCommand(
		cmd.BuildCmd,
		cmd.StartCmd,
		cmd.StopCmd,
		cmd.UpdateCmd,
		cmd.LogCmd,
		cmd.TestCmd,
	)

	err := cmd.Init()
	if err != nil {
		log.Fatalf("Encountered init error %v", err)
	}

	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
