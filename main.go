package main

import (
	"fmt"

	docker "github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const IMAGE_NAME string = "tomiyou/crypto-arbitrage:latest"

var rootCmd = &cobra.Command{
	Use:   "crypto-arbitrage",
	Short: "Crypto-arbitrage bot written in go.",
}

var client *docker.Client

func init() {
	var err error

	// init cobra commands
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(cleanCmd)

	// init docker api
	client, err = docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
