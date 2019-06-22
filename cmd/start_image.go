package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start crypto-arbitrage bot.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		////////////////////// make the user choose the config that is used as base /////////////////////////////////
		// first we read the folder contents
		files, err := ioutil.ReadDir("config")
		if err != nil {
			return
		}

		// then we eliminate all non configs
		count := 0
		configs := make([]string, len(files))
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".yaml") {
				configs[count] = f.Name()
				count += 1
			}
		}
		configs = configs[:count]

		// now we let the user choose config
		prompt := promptui.Select{
			Label: "Select Config",
			Items: configs,
		}
		_, config, err := prompt.Run()
		if err != nil {
			return
		}

		////////////////////// config was chosen, now build the image with config ///////////////////////////////////
		// remove the .yaml suffix for image name
		imageName := "cryptobot_" + strings.ToLower(config[:len(config)-5])
		configPath := "config/" + config
		ctx := context.Background()

		// create tar file for docker image build
		buildContext, err := createTarFile(botConfig.ArbitrageSrcPath+"/docker/docker-run.Dockerfile", "keys", "config")
		defer buildContext.Close()
		if err != nil {
			return
		}

		// image options
		buildOptions := types.ImageBuildOptions{
			SuppressOutput: false,
			// Remove:         true,
			// ForceRemove:    true,
			// PullParent:     true,
			Tags:       []string{imageName},
			Dockerfile: "docker-run.Dockerfile",
			BuildArgs:  map[string]*string{"configPath": &configPath},
		}

		// build the image
		buildResponse, err := dockerClient.ImageBuild(ctx, buildContext, buildOptions)
		if err != nil {
			return
		}

		err = outputStream(buildResponse.Body)
		if err != nil {
			return
		}

		////////////////////// stop and remove any previous container that uses this imageName //////////////////////
		options := filters.NewArgs()
		options.Add("label", "config_name="+config)

		// first we get the running containers
		containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
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
		csv_folder := filepath.Dir(ex) + "/csv"

		// create the container
		createContResp, err := dockerClient.ContainerCreate(ctx, &container.Config{
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
					Source: csv_folder,
					Target: "/mounted",
				},
			},
		}, nil, "")
		if err != nil {
			return
		}

		fmt.Println("Created container with ID:", createContResp.ID)

		// run the created container
		if err = dockerClient.ContainerStart(ctx, createContResp.ID, types.ContainerStartOptions{}); err != nil {
			return
		}

		return removeDanglingImages()
	},
}
