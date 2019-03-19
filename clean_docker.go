package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean crypto-arbitrage bot docker files.",
	RunE:  cleanDockerCmd,
}

func cleanDockerCmd(cmd *cobra.Command, args []string) (err error) {
	////////////////////// make the user choose the config that is used as base /////////////////////////////////
	ctx := context.Background()
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})

	for _, container := range containers {
		if container.State != "running" {
			err = stopAndRemoveContainer(container.ID)
			if err != nil {
				return
			}
		}
	}

	////////////////////// the container was chosen, time to stop ///////////////////////////////////////////////
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.HasPrefix(tag, "cryptobot_") {
				_, err = client.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{})
				if err == nil {
					fmt.Println("Removed image with tag:", tag)
				}

				break
			}
		}
	}

	return removeDanglingImages()
}
