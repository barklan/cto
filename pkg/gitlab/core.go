package gitlab

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const baseAPIURL = "https://gitlab.com/api/v4"

func request(gitlabMethod string) ([]byte, error) {
	baseURL := fmt.Sprintf("%s%s%s", baseAPIURL, "/projects/", os.Getenv("GITLAB_PROJECT_ID"))
	req, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/%s", baseURL, gitlabMethod), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GITLAB_API_TOKEN")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("GitLab client get failed: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}
