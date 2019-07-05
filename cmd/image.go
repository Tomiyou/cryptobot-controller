package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func buildReleaseImage() (err error) {
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
	buildResponse, err := dockerClient.docker.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return
	}

	// print the response
	return displayDockerStream(buildResponse.Body)
}

func pushToDockerHub() (err error) {
	// get the credentials from the keyfile
	err = getDockerHubCredentials()
	if err != nil {
		return
	}

	// push the image to docker hub
	pushResponse, err := dockerClient.docker.ImagePush(
		context.Background(),
		"docker.io/"+botConfig.RemoteImageName,
		types.ImagePushOptions{
			RegistryAuth: dockerClient.Auth64,
		},
	)
	if err != nil {
		return
	}

	// print the response
	return displayDockerStream(pushResponse)
}

func removeDanglingImages() (err error) {
	filters := filters.NewArgs()
	filters.Add("dangling", "true")

	// first get dangling images
	images, err := dockerClient.docker.ImageList(context.Background(), types.ImageListOptions{Filters: filters})
	if err != nil {
		return
	}

	// now remove the images and their children
	fmt.Println("Removing dangling images:")
	for _, image := range images {
		removedImages, err := dockerClient.docker.ImageRemove(context.Background(), image.ID[7:], types.ImageRemoveOptions{
			PruneChildren: true,
		})
		if err != nil {
			return err
		}

		// quickly loop through all the removed images and print their ids
		for _, removedImage := range removedImages {
			fmt.Println(removedImage.Deleted)
		}
	}

	return
}
