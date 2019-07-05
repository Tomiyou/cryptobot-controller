package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		////////////////////// make the user choose the config that is used as base /////////////////////////////////
		config, err := chooseConfigFile()

		////////////////////// config was chosen, now build the image with config ///////////////////////////////////
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
			PullParent:     true,
			Tags:           []string{imageName},
			Dockerfile:     "docker-run.Dockerfile",
			BuildArgs: map[string]*string{
				"configPath":    &configPath,
				"orgConfigName": &configName,
			},
		}

		// build the image
		buildResponse, err := dockerClient.docker.ImageBuild(ctx, buildContext, buildOptions)
		if err != nil {
			return
		}

		err = displayDockerStream(buildResponse.Body)
		if err != nil {
			return
		}

		////////////////////// stop and remove any previous container that uses this imageName //////////////////////
		options := filters.NewArgs()
		options.Add("label", "config_name="+config)

		// first we get the running containers
		containers, err := dockerClient.docker.ContainerList(ctx, types.ContainerListOptions{
			All:     true,
			Filters: options,
		})
		if err != nil {
			return
		}

		for _, container := range containers {
			// stop the container
			err = stopAndRemoveContainer(container.ID)
			if err != nil {
				return err
			}

			fmt.Println("Stopped and removed container with id:", container.ID)
		}

		////////////////////// now run the created image ////////////////////////////////////////////////////////////
		// get the cwd used for mounting csv folder
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		csvFolder := filepath.Dir(ex) + "/csv"
		_ = os.Mkdir(csvFolder, 0777)

		// create the container
		createContResp, err := dockerClient.docker.ContainerCreate(ctx, &container.Config{
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
		}, nil, "")
		if err != nil {
			return
		}

		fmt.Println("Created container with ID:", createContResp.ID)

		// run the created container
		if err = dockerClient.docker.ContainerStart(ctx, createContResp.ID, types.ContainerStartOptions{}); err != nil {
			return
		}

		return removeDanglingImages()
	},
}
