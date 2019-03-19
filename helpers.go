package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Tomiyou/jsonLoader"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/mholt/archiver"
)

func stopAndRemoveContainer(id string) (err error) {
	ctx := context.Background()
	// stop the container
	if err = dockerClient.ContainerStop(ctx, id, nil); err != nil {
		return
	}
	// remove the container
	if err = dockerClient.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return
	}

	return
}

func removeDanglingImages() (err error) {
	filters := filters.NewArgs()
	filters.Add("dangling", "true")

	// first get dangling images
	images, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{Filters: filters})
	if err != nil {
		return
	}

	// now remove the images and their children
	fmt.Println("Removed:")
	for _, image := range images {
		removedImages, err := dockerClient.ImageRemove(context.Background(), image.ID[7:], types.ImageRemoveOptions{
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

func outputStream(body io.ReadCloser) (err error) {
	defer body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(body, os.Stderr, termFd, isTerm, nil)

	return
}

// !!! IMPORTANT
// We don't close the file, because we depend on the user to do so
func createTarFile(inputs ...string) (file *os.File, err error) {
	// first create the archive
	tar := archiver.Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	err = tar.Archive(inputs, config.TemporaryTarPath)
	if err != nil {
		return
	}

	// now open the tar file using the reader interface (archivex.TarFile has writer interface)
	if file, err = os.Open(config.TemporaryTarPath); err != nil {
		return
	}

	return
}

func extractTarFile(destination string) (err error) {
	tar := archiver.Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	err = tar.Unarchive(config.TemporaryTarPath, destination)
	return
}

func getDockerHubCredentials() (result string, err error) {
	authConfig := types.AuthConfig{}
	if err = jsonLoader.LoadJSON("keys/docker-auth.key", &authConfig); err != nil {
		return
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return
	}
	result = base64.URLEncoding.EncodeToString(encodedJSON)

	return
}
