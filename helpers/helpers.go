package helpers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"bitbucket.org/tomihrib/crypto-arbitrage/helpers"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run latest test function.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// run the created container
		if err = client.ContainerStart(context.Background(), args[0], types.ContainerStartOptions{}); err != nil {
			return
		}

		return
	},
}

const TMP_TAR_PATH string = "/tmp/cryptobot_temporary.tar"

func stopAndRemoveContainer(id string) (err error) {
	ctx := context.Background()
	// stop the container
	if err = client.ContainerStop(ctx, id, nil); err != nil {
		return
	}
	// remove the container
	if err = client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return
	}

	return
}

func removeDanglingImages() (err error) {
	filters := filters.NewArgs()
	filters.Add("dangling", "true")

	// first get dangling images
	images, err := client.ImageList(context.Background(), types.ImageListOptions{Filters: filters})
	if err != nil {
		return
	}

	// now remove the images and their children
	fmt.Println("Removed:")
	for _, image := range images {
		removedImages, err := client.ImageRemove(context.Background(), image.ID[7:], types.ImageRemoveOptions{
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

	return
}

func outputStream(body io.ReadCloser) (err error) {
	defer body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(body, os.Stderr, termFd, isTerm, nil)

	return
}

// !!! IMPORTANT
// We don't close the file, because we depend on the user to do so
func createTarFile(inputs ...string) (file *os.File, err error) {
	// first create the archive
	tar := archiver.Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	err = tar.Archive(inputs, TMP_TAR_PATH)
	if err != nil {
		return
	}

	// now open the tar file using the reader interface (archivex.TarFile has writer interface)
	if file, err = os.Open(TMP_TAR_PATH); err != nil {
		return
	}

	return
}

func extractTarFile(destination string) (err error) {
	tar := archiver.Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	err = tar.Unarchive(TMP_TAR_PATH, destination)
	return
}

func getDockerHubCredentials() (result string, err error) {
	authConfig := types.AuthConfig{}
	if err = helpers.LoadJSON("keys/docker-auth.key", &authConfig); err != nil {
		return
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return
	}
	result = base64.URLEncoding.EncodeToString(encodedJSON)

	return
}
