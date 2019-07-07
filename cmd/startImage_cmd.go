package cmd

import (
	"bufio"
	"context"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// first handle all the needed user input
		config, err := chooseConfigFile()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter container name: ")
		containerName, _ := reader.ReadString('\n')
		containerName = containerName[:len(containerName)-1]
		fmt.Println("ime:", containerName)

		// remove the .yaml suffix for image name
		configName := strings.ToLower(config[:strings.LastIndex(config, ".")])
		configPath := "config/" + config
		imageName := "cryptobot_" + configName
		ctx := context.Background()

		// create tar file for docker image build
		buildContext, err := createTarFile("dockerfiles/arbitrage/docker-run.Dockerfile", "keys", "config")
		defer buildContext.Close()
		if err != nil {
			return
		}

		// image options
		buildOptions := types.ImageBuildOptions{
			SuppressOutput: false,
			Remove:         true,
			ForceRemove:    true,
			Tags:           []string{imageName},
			Dockerfile:     "docker-run.Dockerfile",
			BuildArgs: map[string]*string{
				"configPath":    &configPath,
				"orgConfigName": &configName,
			},
		}

		// build the image
		buildResponse, err := client.api.ImageBuild(ctx, buildContext, buildOptions)
		if err != nil {
			return
		}

		err = displayDockerStream(buildResponse.Body)
		if err != nil {
			return
		}

		// get the cwd used for mounting csv folder
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		csvFolder := filepath.Dir(ex) + "/csv"
		_ = os.Mkdir(csvFolder, 0777)

		// create the container
		createContResp, err := client.api.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			// Tty:    true,
			Labels: map[string]string{
				"config_name": config,
			},
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
			return
		}

		fmt.Println("Created container with name:", containerName, "and ID:", createContResp.ID)

		// run the created container
		if err = client.api.ContainerStart(ctx, createContResp.ID, types.ContainerStartOptions{}); err != nil {
			return
		}

		return removeDanglingImages()
	},
}
