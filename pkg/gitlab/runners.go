package gitlab

import (
	"encoding/json"
	"fmt"
	"log"
)

type Runner struct {
	Active      bool   `json:"active"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
	IPAddress   string `json:"ip_address"`
	IsShared    bool   `json:"is_shared"`
	Name        string `json:"name"`
	Online      bool   `json:"online"`
	RunnerType  string `json:"runner_type"`
	Status      string `json:"status"`
}

func GetActiveGroupRunners(gitlabProjectId, gitlabToken string) ([]Runner, error) {
	runners := make([]Runner, 1)
	dump, err := request(
		gitlabProjectId,
		gitlabToken,
		fmt.Sprintf("runners?status=active&type=group_type"),
	)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(dump))
	err = json.Unmarshal(dump, &runners)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return runners, nil
}
