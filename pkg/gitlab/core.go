package gitlab

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const baseAPIURL = "https://gitlab.com/api/v4"

func request(projectId, token, gitlabMethod string) ([]byte, error) {
	baseURL := fmt.Sprintf("%s%s%s", baseAPIURL, "/projects/", projectId)
	req, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/%s", baseURL, gitlabMethod), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("GitLab client get failed: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}
