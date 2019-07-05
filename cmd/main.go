package cmd

import (
	"github.com/Tomiyou/jsonLoader"
	docker "github.com/docker/docker/client"
)

var dockerClient struct {
	docker *docker.Client
	Auth64 string
}

var botConfig struct {
	RemoteImageName      string `json:"remoteImageName"`
	TemporaryTarPath     string `json:"temporaryTarPath"`
	EncryptedSecretsPath string `json:"encryptedSecretsPath"`
	ArbitrageSrcPath     string `json:"ArbitrageSrcPath"`
}

func Initialize() (err error) {
	// load config
	err = jsonLoader.LoadJSON("settings.json", &botConfig)
	if err != nil {
		panic(err)
	}

	// add / to arbitrage src folder
	if botConfig.ArbitrageSrcPath[len(botConfig.ArbitrageSrcPath)-1] != '/' {
		botConfig.ArbitrageSrcPath += "/"
	}

	// init docker api
	client, err := docker.NewEnvClient()
	if err != nil {
		return
	}

	dockerClient.docker = client

	return
}
