package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Run runs arbitrary container with specified command and
// returns string containing stdout ignoring stderr
func Run(baseImage, tag string, command []string) string {
	image := fmt.Sprintf("%s:%s", baseImage, tag)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	authConfig := types.AuthConfig{
		Username: os.Getenv("GITLAB_TOKEN_USERNAME"),
		Password: os.Getenv("GITLAB_TOKEN_PASSWORD"),
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	// encodedJSON, _ := json.Marshal(map[string]string{
	// 	"username": os.Getenv("GITLAB_TOKEN_USERNAME"),
	// 	"password": os.Getenv("GITLAB_TOKEN_PASSWORD"),
	// })
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   command,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	buf := new(bytes.Buffer)
	stdcopy.StdCopy(buf, os.Stderr, out)
	content := buf.String()

	return content
}
