package main

import (
	"fmt"

	"github.com/Tomiyou/jsonLoader"
	docker "github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crypto-arbitrage",
	Short: "Crypto-arbitrage bot written in go.",
}

var dockerClient *docker.Client

var botConfig struct {
	RemoteImageName      string `json:"remoteImageName"`
	TemporaryTarPath     string `json:"temporaryTarPath"`
	EncryptedSecretsPath string `json:"encryptedSecretsPath"`
	CryptobotSrcPath     string `json:"cryptobotSrcPath"`
}

func init() {
	var err error

	// load config
	err = jsonLoader.LoadJSON("config.json", &botConfig)
	if err != nil {
		panic(err)
	}

	// init cobra commands
	rootCmd.AddCommand(buildCmd, startCmd, stopCmd, logCmd, updateCmd, encryptCmd, decryptCmd, cleanCmd)

	// init docker api
	dockerClient, err = docker.NewEnvClient()
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
