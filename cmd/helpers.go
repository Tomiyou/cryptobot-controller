package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Tomiyou/jsonLoader"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/manifoldco/promptui"
	"github.com/mholt/archiver"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

func stopAndRemoveContainer(id string) (err error) {
	ctx := context.Background()
	// stop the container
	if err = dockerClient.docker.ContainerStop(ctx, id, nil); err != nil {
		return
	}
	// remove the container
	if err = dockerClient.docker.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return
	}

	return
}

// IMPORTANT: .close() method needs to be called manually; also '/' at the end of the name is for folders
func createTarFile(inputs ...string) (file *os.File, err error) {
	// create the archive
	tar := archiver.Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	err = tar.Archive(inputs, botConfig.TemporaryTarPath)
	if err != nil {
		return
	}

	// open the created tar file and return the interface
	return os.Open(botConfig.TemporaryTarPath)
}

// read docker credentials from file
func getDockerHubCredentials() (err error) {
	if dockerClient.Auth64 != "" {
		return
	}

	authConfig := types.AuthConfig{}
	if err = jsonLoader.LoadJSON("keys/docker-auth.key", &authConfig); err != nil {
		return
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return
	}
	dockerClient.Auth64 = base64.URLEncoding.EncodeToString(encodedJSON)

	return
}

// Pass the stream from a io.ReadCloser interface
func displayDockerStream(body io.ReadCloser) (err error) {
	defer body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(body, os.Stderr, termFd, isTerm, nil)

	return
}

func chooseContainer() (container types.Container, err error) {
	ctx := context.Background()
	containers, err := dockerClient.docker.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})

	if len(containers) == 0 {
		return types.Container{}, fmt.Errorf("No containers present.")
	}

	// now we let the user choose container
	imageNames := make([]string, len(containers))
	for i, _ := range containers {
		imageNames[i] = containers[i].Status + " : " + containers[i].Image
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

func chooseConfigFile() (config string, err error) {
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
	_, config, err = prompt.Run()

	return
}
