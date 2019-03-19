package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
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
	reader, err := client.ImagePull(context.Background(), "docker.io/"+IMAGE_NAME, types.ImagePullOptions{
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