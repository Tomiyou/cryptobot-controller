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
		////////////////////// create tar file for docker image build ///////////////////////////////////////////////
		buildContext, err := createTarFile(
			botConfig.ArbitrageSrcPath,
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
			Tags:           []string{botConfig.RemoteImageName},
			Dockerfile:     "docker-build.Dockerfile",
		}

		// build the image
		ctx := context.Background()
		buildResponse, err := dockerClient.ImageBuild(ctx, buildContext, buildOptions)
		if err != nil {
			return
		}

		err = outputStream(buildResponse.Body)
		if err != nil {
			return
		}

		////////////////////// now push the created image to docker hub for remote access ///////////////////////////
		// get the credentials from the keyfile
		auth64, err := getDockerHubCredentials()
		if err != nil {
			return
		}

		// push the image to docker hub
		pushResponse, err := dockerClient.ImagePush(context.Background(), "docker.io/"+botConfig.RemoteImageName, types.ImagePushOptions{
			RegistryAuth: auth64,
		})
		if err != nil {
			return
		}

		err = outputStream(pushResponse)
		if err != nil {
			return
		}

		////////////////////// remove any leftover images ///////////////////////////////////////////////////////////
		return removeDanglingImages()
	},
}
