package cmd

import (
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a crypto-arbitrage bot image.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// create tar file for docker image build
		err = buildReleaseImage()
		if err != nil {
			return
		}

		// now push the created image to docker hub for remote access
		err = pushToDockerHub()
		if err != nil {
			return
		}

		// remove any leftover images
		return removeDanglingImages()
	},
}
