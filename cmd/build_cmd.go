package cmd

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a crypto-arbitrage bot image.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// create tar file for docker image build
		buildContext, err := createTarFile(config.CryptobotSource, "dockerfiles/build.Dockerfile")
		defer buildContext.Close()
		if err != nil {
			return err
		}

		// image options
		buildOptions := types.ImageBuildOptions{
			SuppressOutput: false,
			Remove:         true,
			ForceRemove:    true,
			PullParent:     true,
			Tags:           []string{config.RemoteImageName},
			Dockerfile:     "build.Dockerfile",
		}

		// build the image
		ctx := context.Background()
		buildResponse, err := client.api.ImageBuild(ctx, buildContext, buildOptions)
		if err != nil {
			return err
		}

		// print the response

		if err := displayDockerStream(buildResponse.Body); err != nil {
			return err
		}

		if err := ensureDockerCredentials(); err != nil {
			return err
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
			return err
		}

		// print the response
		if err := displayDockerStream(pushResponse); err != nil {
			return err
		}

		// remove any leftover images
		return removeDanglingImages()
	},
}
