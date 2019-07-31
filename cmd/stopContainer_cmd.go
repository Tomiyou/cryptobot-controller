package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// make the user choose the config that is used as base
		container, err := chooseContainer()
		if err != nil {
			return err
		}

		// the container was chosen, time to stop it
		ctx := context.Background()
		if err := client.api.ContainerStop(ctx, container.ID, nil); err != nil {
			return err
		}

		// remove the container
		if err := client.api.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}

		fmt.Println("Stopped and removed container with name:", container.Names[0], "and ID:", container.ID)

		// Remove the image belonging to the container too
		_, err = client.api.ImageRemove(context.Background(), container.ImageID, types.ImageRemoveOptions{
			PruneChildren: true,
		})
		if err != nil {
			return err
		}

		fmt.Println("Removed image with ID:", container.ImageID)

		return nil
	},
}
