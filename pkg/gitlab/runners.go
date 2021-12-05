package gitlab

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

// func GetActiveGroupRunners() ([]Runner, error) {
// 	runners := make([]Runner, 1)
// 	dump, _ := request(fmt.Sprintf("runners?status=active&type=group_type"))
// 	fmt.Println(string(dump))
// 	err := json.Unmarshal(dump, &runners)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return runners, nil
// }
