package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update crypto-arbitrage bot.",
	RunE:  updateImageCmd,
}

func updateImageCmd(cmd *cobra.Command, args []string) (err error) {
	////////////////////// pull the image from the docker hub ///////////////////////////////////////////////////
	// get the credentials from the keyfile
	auth64, err := getDockerHubCredentials()
	if err != nil {
		return
	}

	// pull the image from the docker hub
	reader, err := dockerClient.ImagePull(context.Background(), "docker.io/"+botConfig.RemoteImageName, types.ImagePullOptions{
		RegistryAuth: auth64,
	})
	if err != nil {
		return
	}

	err = outputStream(reader)
	if err != nil {
		return
	}

	return
}
