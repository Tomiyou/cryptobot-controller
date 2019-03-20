package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop crypto-arbitrage bot.",
	RunE:  stopImageCmd,
}

func stopImageCmd(cmd *cobra.Command, args []string) (err error) {
	////////////////////// make the user choose the config that is used as base /////////////////////////////////
	container, err := chooseContainer()
	if err != nil {
		return
	}

	////////////////////// the container was chosen, time to stop ///////////////////////////////////////////////
	err = stopAndRemoveContainer(container.ID)
	if err != nil {
		return
	}

	fmt.Println("Stopped and removed container with id:", container.ID)

	return
}

func chooseContainer() (container types.Container, err error) {
	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})

	if len(containers) == 0 {
		return types.Container{}, fmt.Errorf("No containers present.")
	}

	// now we let the user choose container
	imageNames := make([]string, len(containers))
	for i, _ := range containers {
		imageNames[i] = containers[i].Image
	}
	prompt := promptui.Select{
		Label: "Select Container",
		Items: imageNames,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return
	}

	container = containers[index]
	return
}
