package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean crypto-arbitrage bot docker files.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		////////////////////// make the user choose the config that is used as base /////////////////////////////////
		ctx := context.Background()
		containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
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
		images, err := dockerClient.ImageList(ctx, types.ImageListOptions{})
		if err != nil {
			return
		}

		for _, image := range images {
			for _, tag := range image.RepoTags {
				if strings.HasPrefix(tag, "cryptobot_") {
					_, err = dockerClient.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{})
					if err == nil {
						fmt.Println("Removed image with tag:", tag)
					}

					break
				}
			}
		}

		return removeDanglingImages()
	},
}
