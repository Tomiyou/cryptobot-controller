package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// get the credentials from the keyfile
		err = getDockerHubCredentials()
		if err != nil {
			return
		}

		// pull the image from the docker hub
		reader, err := client.api.ImagePull(
			context.Background(),
			"docker.io/"+config.RemoteImageName,
			types.ImagePullOptions{
				RegistryAuth: client.Auth64,
			},
		)
		if err != nil {
			return
		}

		err = displayDockerStream(reader)
		if err != nil {
			return
		}

		return
	},
}
