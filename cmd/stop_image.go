package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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
	},
}

