package tempnginx

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/barklan/cto/pkg/manager/exec"
	"github.com/barklan/cto/pkg/storage"
	"github.com/foomo/htpasswd"
)

func buildNginx(
	data *storage.Data,
	projectName,
	basicAuthUsername,
	basicAuthPassword string,
) error {
	builderPath := data.CreateMediaDirIfNotExists("tempnginx")

	dockerFilePath := builderPath + "/Dockerfile"
	dockerFileStr := `FROM nginx:1.19.0-alpine
ARG BUILDKIT_INLINE_CACHE=1
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d
COPY htpasswd /etc/nginx/htpasswd
`

	if err := os.WriteFile(dockerFilePath, []byte(dockerFileStr), 0777); err != nil {
		return err
	}

	nginxFilePath := builderPath + "/nginx.conf"
	nginxFileStr := fmt.Sprintf(`server {

	auth_basic           "Restricted Area";
	auth_basic_user_file /etc/nginx/htpasswd;

    listen 9090;
    client_max_body_size 5M;

    location / {
        alias /home/app/media/%s;
        autoindex on;
    }
}
`, "")

	if err := os.WriteFile(nginxFilePath, []byte(nginxFileStr), 0777); err != nil {
		return err
	}

	htpasswdPath := builderPath + "/htpasswd"
	file, err := os.Create(htpasswdPath)
	if err != nil {
		log.Panic(err)
	}
	file.Close()

	err = htpasswd.SetPassword(htpasswdPath, basicAuthUsername, basicAuthPassword, htpasswd.HashBCrypt)
	if err != nil {
		log.Panic(err)
	}

	cmd := []string{"docker", "build", "-t", fmt.Sprintf("nginx:%s", projectName), builderPath}
	_, err = exec.ExecNoShell(cmd)
	return err
}

func TemporaryNginx(
	data *storage.Data,
	projectName string,
	addressChan chan string,
	minutes int,
	basicAuthUsername,
	basicAuthPassword string,
) {
	data.CreateMediaDirIfNotExists(projectName)

	if err := buildNginx(data, projectName, basicAuthUsername, basicAuthPassword); err != nil {
		log.Println(err)
		data.PSend(projectName, "Failed to build temporary nginx server.")
		return
	}

	// TODO randomize port and name to allow for multiple nginx servers running simultaneously
	containerName := "tempnginx"
	port := "9090"

	cmd := []string{
		"docker",
		"run",
		"--name",
		containerName,
		"--rm",
		"-d",
		"-v", "cto-media:/home/app/media",
		"-p", fmt.Sprintf("%s:%s", port, port),
		fmt.Sprintf("nginx:%s", projectName),
	}

	_, err := exec.ExecNoShell(cmd)
	if err != nil {
		log.Println(err)
		data.PSend(projectName, "Failed to start temporary nginx server.")
		return
	}

	// TODO this should not be hardcoded.
	hostname := "ctopanel.com"
	addressChan <- fmt.Sprintf("%s:%s", hostname, port)

	time.Sleep(time.Duration(minutes) * time.Minute)

	_, err = exec.ExecNoShell(
		[]string{"docker", "stop", containerName},
	)
	if err != nil {
		data.PSend(projectName, "Failed to stop temporary nginx server. Please intervene.")
		return
	}

	data.PSend(projectName, "Temporary nginx server stopped as planned.")
}
