package main

import (
	"fmt"

	"github.com/Tomiyou/cryptobot-controller/cmd"

	"github.com/Tomiyou/jsonLoader"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crypto-arbitrage",
	Short: "Crypto-arbitrage bot written in go.",
}

func init() {
	// load config
	err := jsonLoader.LoadJSON("settings.json", &Config)
	if err != nil {
		panic(err)
	}

	// add / to arbitrage src folder
	if Config.ArbitrageSrcPath[len(Config.ArbitrageSrcPath)-1] != '/' {
		Config.ArbitrageSrcPath += "/"
	}

	// init docker api
	DockerClient, err = docker.NewEnvClient()
	if err != nil {
		return
	}

	// init cobra commands
	rootCmd.AddCommand(
		cmd.BuildCmd,
		cmd.StartCmd,
		cmd.StopCmd,
		cmd.UpdateCmd,
		cmd.CleanCmd,
		cmd.LogCmd,
		cmd.TestCmd,
	)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
