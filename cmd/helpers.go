package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Tomiyou/jsonLoader"
	"github.com/manifoldco/promptui"
	"github.com/mholt/archiver"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

// Let the user choose containers from a list
func chooseContainer() (types.Container, error) {
	ctx := context.Background()
	containers, err := client.api.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})

	if len(containers) == 0 {
		return types.Container{}, fmt.Errorf("No containers present.")
	}

	// now we let the user choose container
	containerNames := make([]string, len(containers))
	for i, container := range containers {
		containerNames[i] = container.Status + " : " + strings.Join(container.Names, ";")
	}
	prompt := promptui.Select{
		Label: "Select Container",
		Items: containerNames,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return types.Container{}, err
	}

	return containers[index], nil
}

// Let the user choose config file
func chooseConfigFile() (string, error) {
	// first we read the folder contents
	files, err := ioutil.ReadDir("config")
	if err != nil {
		return "", err
	}

	// then we eliminate all non configs
	configs := make([]string, len(files))[:0]
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			configs = append(configs, f.Name())
		}
	}

	// now we let the user choose config
	prompt := promptui.Select{
		Label: "Select Config",
		Items: configs,
	}
	_, config, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return config, nil
}

// IMPORTANT: .close() method needs to be called manually; also '/' at the end of the name is for folders
func createTarFile(inputs ...string) (file *os.File, err error) {
	// create the archive
	tar := archiver.Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	err = tar.Archive(inputs, config.TemporaryTarPath)
	if err != nil {
		return
	}

	// open the created tar file and return the interface
	return os.Open(config.TemporaryTarPath)
}

// read docker credentials from file
func getDockerHubCredentials() error {
	if client.Auth64 != "" {
		return nil
	}

	authConfig := types.AuthConfig{}
	if err := jsonLoader.LoadJSON("keys/docker-auth.key", &authConfig); err != nil {
		return err
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	client.Auth64 = base64.URLEncoding.EncodeToString(encodedJSON)

	return nil
}

// Print the stream from a io.ReadCloser interface to stdout
func displayDockerStream(body io.ReadCloser) error {
	defer body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stdout)
	return jsonmessage.DisplayJSONMessagesStream(body, os.Stdout, termFd, isTerm, nil)
}

func removeDanglingImages() error {
	filters := filters.NewArgs()
	filters.Add("dangling", "true")

	// first get dangling images
	images, err := client.api.ImageList(context.Background(), types.ImageListOptions{Filters: filters})
	if err != nil {
		return err
	}

	// now remove the images and their children
	fmt.Println("Removing dangling images:")
	for _, image := range images {
		removedImages, err := client.api.ImageRemove(context.Background(), image.ID[7:], types.ImageRemoveOptions{
			PruneChildren: true,
		})
		if err != nil {
			return err
		}

		// quickly loop through all the removed images and print their ids
		for _, removedImage := range removedImages {
			fmt.Println(removedImage.Deleted)
		}
	}

	return nil
}
