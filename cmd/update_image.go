package cmd

import (
	"context"

	"github.com/Tomiyou/cryptobot-controller/api"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		////////////////////// pull the image from the docker hub ///////////////////////////////////////////////////
		// get the credentials from the keyfile
		auth64, err := api.GetDockerHubCredentials()
		if err != nil {
			return
		}

		// pull the image from the docker hub
		reader, err := DockerClient.ImagePull(context.Background(), "docker.io/"+Config.RemoteImageName, types.ImagePullOptions{
			RegistryAuth: auth64,
		})
		if err != nil {
			return
		}

		err = api.DisplayDockerStream(reader)
		if err != nil {
			return
		}

		return
	},
}
