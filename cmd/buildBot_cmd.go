package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a crypto-arbitrage bot image.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// create tar file for docker image build
		buildContext, err := createTarFile(
			config.ArbitrageSrcPath,
			"dockerfiles/arbitrage/docker-build.Dockerfile")
		defer buildContext.Close()
		if err != nil {
			return
		}

		// image options
		buildOptions := types.ImageBuildOptions{
			SuppressOutput: false,
			Remove:         true,
			ForceRemove:    true,
			PullParent:     true,
			Tags:           []string{config.RemoteImageName},
			Dockerfile:     "docker-build.Dockerfile",
		}

		// build the image
		ctx := context.Background()
		buildResponse, err := client.api.ImageBuild(ctx, buildContext, buildOptions)
		if err != nil {
			return
		}

		// print the response
		err = displayDockerStream(buildResponse.Body)
		if err != nil {
			return
		}

		// now push the created image to docker hub for remote access
		// get the credentials from the keyfile
		err = getDockerHubCredentials()
		if err != nil {
			return
		}

		// push the image to docker hub
		pushResponse, err := client.api.ImagePush(
			context.Background(),
			"docker.io/"+config.RemoteImageName,
			types.ImagePushOptions{
				RegistryAuth: client.Auth64,
			},
		)
		if err != nil {
			return
		}

		// print the response
		err = displayDockerStream(pushResponse)
		if err != nil {
			return
		}

		// remove any leftover images
		return removeDanglingImages()
	},
}
