package querying

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/barklan/cto/pkg/storage"
)

type RequestQuery struct {
	ProjectName string
	Env         string
	Service     string
	Date        string
	TimeQuery   string
}

func (rq RequestQuery) BeaconToSeek() string {
	timeQueryBeacon := TimeQueryBeaconToSeek(rq.TimeQuery)
	beacon := strings.Join([]string{
		rq.ProjectName,
		rq.Env,
		rq.Service,
		rq.Date,
		timeQueryBeacon,
	}, " ")
	return beacon
}

func (rq RequestQuery) ValidPrefix() string {
	return strings.Join([]string{
		rq.ProjectName,
		rq.Env,
		rq.Service,
		rq.Date,
		rq.TimeQuery,
	}, " ")
}

func TimeQueryBeaconToSeek(timeQuery string) string {
	if timeQuery == "" {
		return "24:00:00"
	}
	padWithLeadingZero := true
	if string(timeQuery[len(timeQuery)-1]) == ":" {
		timeQuery = timeQuery[:len(timeQuery)-1]
	} else if string(timeQuery[len(timeQuery)-2]) == ":" {
		padWithLeadingZero = false
	}

	timeQuerySplit := strings.Split(timeQuery, ":")
	last := timeQuerySplit[len(timeQuerySplit)-1]
	lastInt, err := strconv.ParseInt(last, 10, 64)
	if err != nil {
		log.Println("failed to parse int in time query:", err)
	}
	lastInt++
	lastIntStr := strconv.FormatInt(lastInt, 10)

	if len(lastIntStr) == 1 && padWithLeadingZero {
		lastIntStr = "0" + lastIntStr
	}
	timeQuerySplit[len(timeQuerySplit)-1] = lastIntStr
	result := strings.Join(timeQuerySplit, ":")
	return result
}

func GetFullEnv(data *storage.Data, project, envQ string) (string, bool) {
	numberOfMatches := 0
	var lastMatch string

	knownEnvs := GetKnownEnvs(data, project)
	for k := range knownEnvs {
		if strings.Contains(k, envQ) {
			numberOfMatches++
			lastMatch = k
		}
	}

	if numberOfMatches != 1 {
		return "", false
	}
	return lastMatch, true
}

func GetFullService(data *storage.Data, project, environment, serviceQ string) (string, bool) {
	numberOfMatches := 0
	var lastMatch string

	knownServices := GetKnownServices(data, project, environment)
	for k := range knownServices {
		if strings.Contains(k, serviceQ) {
			numberOfMatches++
			lastMatch = k
		}
	}

	if numberOfMatches != 1 {
		return "", false
	}
	return lastMatch, true
}

func PlaceQuery(w http.ResponseWriter, r *http.Request, data *storage.Data, queueChan chan QueryJob) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	rawQuery := r.URL.Query()

	powerTokensQ := rawQuery.Get("powertoken")
	powerTokens := strings.Fields(powerTokensQ)

	tokenQ := rawQuery.Get("token")
	projectName, statusCode, ok := authorize(data, tokenQ)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}

	query := rawQuery.Get("query")
	log.Printf("Requested query %q", query)

	querySet := strings.Fields(query)
	queryLen := len(querySet)

	if queryLen < 3 || queryLen > 5 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	environmentQ := querySet[0]
	environment, ok := GetFullEnv(data, projectName, environmentQ)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	serviceQ := querySet[1]
	service, ok := GetFullService(data, projectName, environment, serviceQ)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO proper validation
	// match, _ := regexp.MatchString(`^(0[1-9]|[12]\d|3[01])$`, querySet[2])
	day := querySet[2]

	if day == "t" {
		dayNumber := time.Now().Day()
		day = strconv.FormatInt(int64(dayNumber), 10)
		if dayNumber < 10 {
			day = "0" + day
		}
	} else if len(day) == 1 {
		day = "0" + day
	}
	yearAndMonthStr := time.Now().Format("2006-01")
	date := fmt.Sprintf("%s-%s", yearAndMonthStr, day)

	var timeQuery string
	switch queryLen {
	case 3:
		timeQuery = ""
	case 4:
		timeQuery = querySet[3]
		if len(timeQuery) == 1 {
			timeQuery = "0" + timeQuery
		} else if string(timeQuery[1]) == ":" {
			timeQuery = "0" + timeQuery
		}
	}

	requestQuery := RequestQuery{
		ProjectName: projectName,
		Env:         environment,
		Service:     service,
		Date:        date,
		TimeQuery:   timeQuery,
	}
	beacon := requestQuery.BeaconToSeek()
	validPrefix := requestQuery.ValidPrefix()

	// var cutOffOptimHour int64 = -1
	// if string(timeQuery[len(timeQuery)-1]) == "m" {
	// 	match, _ := regexp.MatchString(`^\d+m$`, timeQuery)
	// 	if match == true {
	// 		minutes, err := strconv.ParseInt(timeQuery[:len(timeQuery)-1], 10, 64)
	// 		if err != nil {
	// 			log.Println("Failed to parse minutes.")
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			return
	// 		}
	// 		now := time.Now()
	// 		if minutes < 180 {
	// 			combinedTimeArr := make([]string, minutes+1)
	// 			for i := int64(0); i <= minutes; i++ {
	// 				time := now.Add(time.Duration(-i) * time.Minute)
	// 				paddedTime := fmt.Sprintf("(%s)", time.Format("15:04"))
	// 				combinedTimeArr[i] = paddedTime
	// 			}
	// 			combinedTimeStr := strings.Join(combinedTimeArr, "|")
	// 			timeQuery = fmt.Sprintf("(%s)", combinedTimeStr)
	// 		} else {
	// 			hours := int64(math.Ceil(float64(minutes) / 60.0))
	// 			combinedTimeArr := make([]string, hours)
	// 			for i := int64(0); i < hours; i++ {
	// 				time := now.Add(time.Duration(-i) * time.Hour)
	// 				paddedTime := fmt.Sprintf("(%s:)", time.Format("15"))
	// 				combinedTimeArr[i] = paddedTime
	// 			}
	// 			combinedTimeStr := strings.Join(combinedTimeArr, "|")
	// 			timeQuery = fmt.Sprintf("(%s)", combinedTimeStr)
	// 		}

	// 		cutOffOptim := now.Add(time.Duration(-minutes) * time.Minute)
	// 		if cutOffOptim.Day() == now.Day() {
	// 			cutOffOptimHour = int64(cutOffOptim.Hour()) - 1
	// 		}
	// 	} else {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}
	// } else {
	// 	requestedHour, err := strconv.ParseInt(strings.Split(querySet[1], ":")[0], 10, 64)
	// 	if err != nil {
	// 		cutOffOptimHour = -1
	// 	} else {
	// 		cutOffOptimHour = requestedHour - 1
	// 	}
	// }

	fieldsQ := rawQuery.Get("fields")
	regexQRaw := rawQuery.Get("regex")
	negateUserRegex := false
	var regexQField, regexQ string
	if regexQRaw != "" {
		if strings.Contains(regexQRaw, "=") {
			regexQSplit := strings.SplitN(regexQRaw, "=", 2)
			regexQField, regexQ = regexQSplit[0], regexQSplit[1]

			if string(regexQField[len(regexQField)-1]) == "!" {
				regexQField = regexQField[:len(regexQField)-1]
				negateUserRegex = true
			}

			_, err := regexp.Compile(regexQ)
			if err != nil {
				log.Printf("failed to compile user regexp; %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	} else {
		regexQField, regexQ = "", ""
	}

	isSimpleQuery := false
	if fieldsQ == "" && regexQ == "" {
		isSimpleQuery = true
	}

	rand.Seed(time.Now().UnixNano())
	queryID := "queryJob-" + strconv.FormatInt(rand.Int63(), 10)

	queryJob := QueryJob{
		ID:              queryID,
		IsSimpleQuery:   isSimpleQuery,
		Beacon:          beacon,
		ValidPrefix:     validPrefix,
		FieldsQ:         fieldsQ,
		RegexQ:          regexQ,
		RegexQField:     regexQField,
		PowerTokens:     powerTokens,
		NegateUserRegex: negateUserRegex,
	}

	queueChan <- queryJob

	respMap := map[string]string{"qid": queryID}
	resp, err := json.Marshal(respMap)
	if err != nil {
		log.Println("failed to marshal response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(resp)
}
