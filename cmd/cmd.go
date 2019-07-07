package cmd

import (
	"github.com/Tomiyou/jsonLoader"

	dockerClient "github.com/docker/docker/client"
)

var client struct {
	api    *dockerClient.Client
	Auth64 string
}

var config struct {
	RemoteImageName      string `json:"remoteImageName"`
	TemporaryTarPath     string `json:"temporaryTarPath"`
	EncryptedSecretsPath string `json:"encryptedSecretsPath"`
	ArbitrageSrcPath     string `json:"ArbitrageSrcPath"`
}

func Init() error {
	// load config
	err := jsonLoader.LoadJSON("settings.json", &config)
	if err != nil {
		return err
	}

	// add / to arbitrage src folder
	if config.ArbitrageSrcPath[len(config.ArbitrageSrcPath)-1] != '/' {
		config.ArbitrageSrcPath += "/"
	}

	// init docker api
	dockerClient, err := dockerClient.NewEnvClient()
	if err != nil {
		return err
	}

	client.api = dockerClient
	return nil
}
