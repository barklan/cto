package traefik

type Service struct {
	LoadBalancer struct {
		PassHostHeader bool `json:"passHostHeader"`
		Servers        []struct {
			URL string `json:"url"`
		} `json:"servers"`
	} `json:"loadBalancer"`
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`
	ServerStatus map[string]string `json:"serverStatus"`
	Status       string            `json:"status"`
	Type         string            `json:"type"`
	UsedBy       []string          `json:"usedBy"`
}

// func GetServices(domain string) ([]Service, error) {
// 	traefikURL := fmt.Sprintf("https://traefik.%s/api/http/services", domain)
// 	req, _ := http.NewRequest(http.MethodGet, traefikURL, nil)
// 	req.SetBasicAuth(
// 		os.Getenv("GWB_TRAEFIK_ADMIN_USERNAME"),
// 		data.GetStr("GWB_TRAEFIK_ADMIN_PASSWORD"),
// 	)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	services := make([]Service, 0)
// 	err = json.Unmarshal(body, &services)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return services, nil
// }

// func CheckIfServiceExists(data *storage.Data, domain, serviceName string) {
// 	services, err := GetServices(data, domain)
// 	for _, service := range services {
// 		if service.Name == serviceName {
// 			log.Println("Traefik service exists.")
// 			return
// 		}
// 	}
// 	data.CSend(fmt.Sprintf("Traefik service %q does not exist. @%s", serviceName, data.SysAdmin))
// }
