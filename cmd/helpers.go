package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
)

func stopAndRemoveContainer(id string) (err error) {
	ctx := context.Background()
	// stop the container
	if err = DockerClient.ContainerStop(ctx, id, nil); err != nil {
		return
	}
	// remove the container
	if err = DockerClient.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return
	}

	return
}
