package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:   "logs",
	Short: "Display logs of a crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// make the user chose the config that is used as base
		container, err := chooseContainer()
		if err != nil {
			return err
		}

		// the container was chosen, show the logs
		ctx := context.Background()
		stream, err := client.api.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
			// Follow:     true,
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			return err
		}
		defer stream.Close()

		// print the logs to console
		reader := bufio.NewReader(stream)
		var line string
		for {
			line, err = reader.ReadString('\n')
			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}

			// first 8 bytes of each line are docker header bytes, need to remove them from print
			if len(line) > 8 {
				fmt.Print(line[8:])
			}
		}
	},
}
