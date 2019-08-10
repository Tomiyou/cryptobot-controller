package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// first handle all the needed user input
		config, err := chooseConfigFile()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter container name: ")
		containerName, _ := reader.ReadString('\n')
		containerName = strings.TrimSuffix(containerName, "\n")

		// remove the .yaml suffix for image name
		configPath := "config/" + config
		imageName := "cryptobot_" + strings.ToLower(strings.TrimSuffix(config, filepath.Ext(config)))
		ctx := context.Background()

		// create tar file for docker image build
		buildContext, err := createTarFile("dockerfiles/run.Dockerfile", "keys", "config")
		defer buildContext.Close()
		if err != nil {
			return err
		}

		// image options
		buildOptions := types.ImageBuildOptions{
			SuppressOutput: false,
			Remove:         true,
			ForceRemove:    true,
			Tags:           []string{imageName},
			Dockerfile:     "run.Dockerfile",
			BuildArgs: map[string]*string{
				"configPath":    &configPath,
				"containerName": &containerName,
			},
		}

		// build the image
		buildResponse, err := client.api.ImageBuild(ctx, buildContext, buildOptions)
		if err != nil {
			return err
		}

		err = displayDockerStream(buildResponse.Body)
		if err != nil {
			return err
		}

		// get the pathToSelf used for mounting csv folder
		csvFolder, err := filepath.Abs("csv")
		if err != nil {
			return err
		}

		err = os.Mkdir(csvFolder, 0777)
		log.Println(err)
		if !os.IsExist(err) {
			return err
		}

		// create the container
		createContainerResp, err := client.api.ContainerCreate(ctx, &container.Config{
			Image: imageName,
		}, &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "on-failure",
			},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: csvFolder,
					Target: "/mounted",
				},
			},
		}, nil, containerName)
		if err != nil {
			return err
		}

		fmt.Println("Created container with name:", containerName, "and ID:", createContainerResp.ID)

		// run the created container
		if err := client.api.ContainerStart(ctx, createContainerResp.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}

		return removeDanglingImages()
	},
}
