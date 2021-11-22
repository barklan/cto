package traefik

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/barklan/cto/pkg/storage"
)

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

func GetServices(data *storage.Data, domain string) []Service {
	traefikURL := fmt.Sprintf("https://traefik.%s/api/http/services", domain)
	req, _ := http.NewRequest(http.MethodGet, traefikURL, nil)
	req.SetBasicAuth(os.Getenv("GWB_TRAEFIK_ADMIN_USERNAME"), data.GetStr("GWB_TRAEFIK_ADMIN_PASSWORD"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		data.CSend("Traefik API request failed on stag.")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		data.CSend("Failed to read body from traefik API response.")
	}

	services := make([]Service, 0)
	err = json.Unmarshal(body, &services)
	if err != nil {
		data.CSend("Failed to unmarshal traefik services.")
	}

	return services
}

func CheckIfServiceExists(data *storage.Data, domain, serviceName string) {
	services := GetServices(data, domain)
	for _, service := range services {
		if service.Name == serviceName {
			log.Println("Traefik service exists.")
			return
		}
	}
	data.CSend(fmt.Sprintf("Traefik service %q does not exist. @%s", serviceName, data.SysAdmin))
}
