package api

import (
	docker "github.com/docker/docker/client"

	"github.com/Tomiyou/jsonLoader"
)

type dockerClient struct {
	API    *docker.Client
	Config config
	Auth64 string
}

type config struct {
	RemoteImageName      string `json:"remoteImageName"`
	TemporaryTarPath     string `json:"temporaryTarPath"`
	EncryptedSecretsPath string `json:"encryptedSecretsPath"`
	ArbitrageSrcPath     string `json:"ArbitrageSrcPath"`
}

func New() (dockerClient, error) {
	// create the client struct
	client := dockerClient{
		Config: config{},
	}

	// load config
	err := jsonLoader.LoadJSON("settings.json", &client.Config)
	if err != nil {
		panic(err)
	}

	// add / to arbitrage src folder
	if client.Config.ArbitrageSrcPath[len(client.Config.ArbitrageSrcPath)-1] != '/' {
		client.Config.ArbitrageSrcPath += "/"
	}

	// init docker api
	api, err := docker.NewEnvClient()
	client.API = api

	return client, err
}
