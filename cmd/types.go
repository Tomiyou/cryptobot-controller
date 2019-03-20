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
	CryptobotSrcPath     string `json:"cryptobotSrcPath"`
}

func Initialize() (err error) {
	// load config
	err = jsonLoader.LoadJSON("settings.json", &botConfig)
	if err != nil {
		panic(err)
	}

	// init docker api
	dockerClient, err = docker.NewEnvClient()
	if err != nil {
		return
	}

	return
}
