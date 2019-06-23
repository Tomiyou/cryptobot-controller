package cmd

import (
	docker "github.com/docker/docker/client"

	"github.com/Tomiyou/jsonLoader"
)

var dockerClient *docker.Client

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
	dockerClient, err = docker.NewEnvClient()
	if err != nil {
		return
	}

	return
}
