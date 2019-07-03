package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/Tomiyou/jsonLoader"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/mholt/archiver"
)

// IMPORTANT: .close() method needs to be called manually; also '/' at the end of the name is for folders
func (this dockerClient) createTarFile(inputs ...string) (file *os.File, err error) {
	// create the archive
	tar := archiver.Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	err = tar.Archive(inputs, this.Config.TemporaryTarPath)
	if err != nil {
		return
	}

	// open the created tar file and return the interface
	return os.Open(this.Config.TemporaryTarPath)
}

func (this dockerClient) getDockerHubCredentials() error {
	authConfig := types.AuthConfig{}
	if err := jsonLoader.LoadJSON("keys/docker-auth.key", &authConfig); err != nil {
		return err
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	this.Auth64 = base64.URLEncoding.EncodeToString(encodedJSON)

	return nil
}

// Pass the stream from a io.ReadCloser interface
func (_ dockerClient) displayDockerStream(body io.ReadCloser) (err error) {
	defer body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(body, os.Stderr, termFd, isTerm, nil)

	return
}
