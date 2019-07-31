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
	CryptobotSource      string `json:"cryptobotSource"`
}

func Init() error {
	// load config
	err := jsonLoader.LoadJSON("settings.json", &config)
	if err != nil {
		return err
	}

	// init docker api
	client.api, err = dockerClient.NewEnvClient()
	if err != nil {
		return err
	}

	return nil
}
