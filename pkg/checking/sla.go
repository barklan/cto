package checking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/barklan/cto/pkg/storage"
)

func getSLAinfo(data *storage.Data, projectName string) string {
	totalRuningTimeKey := fmt.Sprintf("%s-totalRunningTime", projectName)
	var totalRunningTime time.Duration
	totalRunningTimeRaw := data.Get(totalRuningTimeKey)
	if string(totalRunningTimeRaw) == "" {
		totalRunningTime = 0
	} else {
		err := json.Unmarshal(totalRunningTimeRaw, &totalRunningTime)
		if err != nil {
			data.CSend("Failed to unmarshal totalRunningTime")
		}
	}

	downTimeKey := fmt.Sprintf("%s-downTime", projectName)
	var totalDownTime time.Duration
	totalDownTimeRaw := data.Get(downTimeKey)
	if string(totalDownTimeRaw) == "" {
		totalDownTime = 0
	} else {
		err := json.Unmarshal(totalDownTimeRaw, &totalDownTime)
		if err != nil {
			data.CSend("Failed to unmarshal totalDownTime")
		}
	}

	invertedSLA := float64(totalDownTime) / float64(totalRunningTime) * 100
	sla := 100.0 - invertedSLA
	return fmt.Sprintf(
		`Uptime: %.10f%% \(total: %s, down: %s\)\. `,
		sla,
		totalRunningTime.Round(time.Second),
		totalDownTime.Round(time.Second),
	)
}

func makeRequest(url string) error {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Not ok.")
	}
	return nil
}

func SLAAggregator(data *storage.Data, projectName string) {
	if url := data.Config.P[projectName].Checks.SLA; url != "" {

		ticker := time.NewTicker(3 * time.Second)

		start := time.Time{}
		zero := time.Time{}
		startCycle := time.Now()

		for {
			err := makeRequest(url)
			if err == nil {
				start = time.Time{}
			} else {
				if start != zero {
					downTime := time.Since(start)
					downTimeKey := fmt.Sprintf("%s-downTime", projectName)
					var totalDownTime time.Duration
					totalDownTimeRaw := data.Get(downTimeKey)
					if string(totalDownTimeRaw) == "" {
						totalDownTime = 0
					} else {
						err := json.Unmarshal(totalDownTimeRaw, &totalDownTime)
						if err != nil {
							data.CSend("Failed to unmarshal totalDownTime")
						}
					}
					totalDownTime += downTime
					data.SetObj(downTimeKey, totalDownTime, -1)
				}
				start = time.Now()
			}

			totalRuningTimeKey := fmt.Sprintf("%s-totalRunningTime", projectName)
			var totalRunningTime time.Duration
			totalRunningTimeRaw := data.Get(totalRuningTimeKey)
			if string(totalRunningTimeRaw) == "" {
				totalRunningTime = 0
			} else {
				err = json.Unmarshal(totalRunningTimeRaw, &totalRunningTime)
				if err != nil {
					data.CSend("Failed to unmarshal totalRunningTime")
				}
			}
			totalRunningTime += time.Since(startCycle)
			startCycle = time.Now()
			data.SetObj(totalRuningTimeKey, totalRunningTime, -1)

			<-ticker.C
		}
	}
}

// main
// go func() {
// 	defer data.CSend("All SLA checks exited.")
// 	wgSLA := new(sync.WaitGroup)
// 	wgSLA.Add(len(data.Config.P))

// 	for projectName := range data.Config.P {
// 		go func(pName string) {
// 			defer func() {
// 				data.CSend(fmt.Sprintf("SLA exited for project %s.", pName))
// 				wgSLA.Done()
// 			}()
// 			checking.SLAAggregator(data, pName)
// 		}(projectName)
// 	}

// 	wgSLA.Wait()
// }()
