package api

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (this dockerClient) buildReleaseImage() (err error) {
	buildContext, err := this.createTarFile(
		this.Config.ArbitrageSrcPath,
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
		Tags:           []string{this.Config.RemoteImageName},
		Dockerfile:     "docker-build.Dockerfile",
	}

	// build the image
	ctx := context.Background()
	buildResponse, err := this.API.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return
	}

	// print the response
	return this.displayDockerStream(buildResponse.Body)
}

func (this dockerClient) pushToDockerHub() (err error) {
	// get the credentials from the keyfile
	err = this.getDockerHubCredentials()
	if err != nil {
		return
	}

	// push the image to docker hub
	pushResponse, err := this.API.ImagePush(
		context.Background(),
		"docker.io/"+this.Config.RemoteImageName,
		types.ImagePushOptions{
			RegistryAuth: this.Auth64,
		},
	)
	if err != nil {
		return
	}

	// print the response
	return this.displayDockerStream(pushResponse)
}

func (this dockerClient) removeDanglingImages() (err error) {
	filters := filters.NewArgs()
	filters.Add("dangling", "true")

	// first get dangling images
	images, err := this.API.ImageList(context.Background(), types.ImageListOptions{Filters: filters})
	if err != nil {
		return
	}

	// now remove the images and their children
	fmt.Println("Removing dangling images:")
	for _, image := range images {
		removedImages, err := this.API.ImageRemove(context.Background(), image.ID[7:], types.ImageRemoveOptions{
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
