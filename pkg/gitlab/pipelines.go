package gitlab

import (
	"encoding/json"
	"fmt"
)

type Pipeline struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	Iid       int64  `json:"iid"`
	ProjectID int64  `json:"project_id"`
	Ref       string `json:"ref"`
	Sha       string `json:"sha"`
	Soure     string `json:"soure"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
	WebURL    string `json:"web_url"`
}

func GetPipelines(branch, status string) ([]Pipeline, error) {
	pipelines := make([]Pipeline, 1)
	dump, _ := request(fmt.Sprintf("pipelines?ref=%s&status=%s", branch, status))
	err := json.Unmarshal(dump, &pipelines)
	if err != nil {
		return nil, err
	}
	return pipelines, nil
}
