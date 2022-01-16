package querying

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/porter"
	"github.com/barklan/cto/pkg/storage"
)

type timeSyntax int

const (
	prefixSyntax timeSyntax = iota
	lastMinutesSyntax
	intervalSyntax
)

type RequestQuery struct {
	ProjectName string
	Env         string
	Service     string
	Date        string
	TimeQuery   string
}

func (rq RequestQuery) BeaconToSeek(syntax timeSyntax, now time.Time) (string, error) {
	var timeQueryBeacon string
	if syntax == lastMinutesSyntax {
		timeQueryBeacon = now.
			Add(1 * time.Minute).
			UTC().
			Format("15:04:05")
	} else {
		timeQueryBeacon = TimeQueryBeaconToSeek(rq.TimeQuery)
	}

	beacon := strings.Join([]string{
		rq.ProjectName,
		rq.Env,
		rq.Service,
		rq.Date,
		timeQueryBeacon,
	}, " ")
	return beacon, nil
}

func (rq RequestQuery) ValidPrefix(syntax timeSyntax) string {
	prefix := strings.Join([]string{
		rq.ProjectName,
		rq.Env,
		rq.Service,
	}, " ")
	if syntax == prefixSyntax {
		prefix += " " + rq.Date + " " + rq.TimeQuery
	}
	return prefix
}

func (rq RequestQuery) NorthStar(
	syntax timeSyntax,
	now time.Time,
) (string, error) {
	// TODO should do it once and not in every function
	if syntax == lastMinutesSyntax {
		minutes, err := strconv.ParseInt(
			rq.TimeQuery[:len(rq.TimeQuery)-1],
			10,
			64,
		)
		if err != nil {
			return "", err
		}
		minutesAgo := now.Add(time.Duration(-minutes) * time.Minute)
		minutesAgoStr := minutesAgo.Format("2006-01-02 15:04:05")
		northStar := strings.Join([]string{
			rq.ProjectName,
			rq.Env,
			rq.Service,
			minutesAgoStr,
		}, " ")

		return northStar, nil
	}
	return "", nil
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
		log.Println("failed to parse int in time query: ", err)
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

func GetFullService(
	data *storage.Data,
	project, environment, serviceQ string,
) (string, bool) {
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

func SetMsgInCache(
	data *storage.Data,
	requestId string,
	status porter.QStatus,
	msg string,
) {
	key := requestId

	valJson, err := json.Marshal(porter.QResp{Msg: msg, Status: status})
	if err != nil {
		log.Panicln("failed to marshal meta message for requested query", err)
	}

	if err := data.Cache.Set(key, valJson, 1*time.Minute); err != nil {
		data.Log.Error("failed to set qmeta in cache", zap.Error(err))
	}
}

func SetResultInCache(
	data *storage.Data,
	requestId string,
	msg string,
	result []map[string]interface{},
) {
	key := requestId

	valJson, err := json.Marshal(porter.QResp{Msg: msg, Status: porter.QDone, Result: result})
	if err != nil {
		log.Panicln("failed to marshal meta message for requested query", err)
	}

	if err := data.Cache.Set(key, valJson, 1*time.Minute); err != nil {
		data.Log.Error("failed to set qmeta in cache", zap.Error(err))
	}
}

func PlaceQuery(
	qr porter.QueryRequest,
	data *storage.Data,
	queueChan chan QueryJob,
) {
	query := qr.QueryText
	log.Printf("Requested query %q", query)

	querySet := strings.Fields(query)
	queryLen := len(querySet)

	if queryLen < 3 || queryLen > 5 {
		SetMsgInCache(
			data,
			qr.RequestID,
			porter.QFailed,
			"Not enough args in main query.",
		)
		return
	}

	environmentQ := querySet[0]
	environment, ok := GetFullEnv(data, qr.ProjectID, environmentQ)
	if !ok {
		SetMsgInCache(
			data,
			qr.RequestID,
			porter.QFailed,
			"No matching environments found.",
		)
		return
	}

	serviceQ := querySet[1]
	service, ok := GetFullService(data, qr.ProjectID, environment, serviceQ)
	if !ok {
		SetMsgInCache(
			data,
			qr.RequestID,
			porter.QFailed,
			"No matching services found.",
		)
		return
	}

	// TODO proper validation
	// match, _ := regexp.MatchString(`^(0[1-9]|[12]\d|3[01])$`, querySet[2])
	day := querySet[2]

	now := time.Now().UTC()

	if day == "t" {
		dayNumber := now.Day()
		day = strconv.FormatInt(int64(dayNumber), 10)
		if dayNumber < 10 {
			day = "0" + day
		}
	} else if len(day) == 1 {
		day = "0" + day
	}
	yearAndMonthStr := now.Format("2006-01")
	date := fmt.Sprintf("%s-%s", yearAndMonthStr, day)

	var timeQuery string
	var tSyntax timeSyntax
	tSyntax = prefixSyntax
	switch queryLen {
	case 3:
		timeQuery = ""
	case 4:
		timeQuery = querySet[3]
		if string(timeQuery[len(timeQuery)-1]) != "m" {
			// TODO this should only be available when t (today) is used as a day
			if len(timeQuery) == 1 {
				timeQuery = "0" + timeQuery
			} else if string(timeQuery[1]) == ":" {
				timeQuery = "0" + timeQuery
			}
		} else {
			tSyntax = lastMinutesSyntax
		}
	}

	requestQuery := RequestQuery{
		ProjectName: qr.ProjectID,
		Env:         environment,
		Service:     service,
		Date:        date,
		TimeQuery:   timeQuery,
	}

	beacon, err := requestQuery.BeaconToSeek(tSyntax, now)
	if err != nil {
		SetMsgInCache(
			data,
			qr.RequestID,
			porter.QFailed,
			"Internal error. Failed to construct beacon.",
		)
		return
	}
	validPrefix := requestQuery.ValidPrefix(tSyntax)

	northStar, err := requestQuery.NorthStar(tSyntax, now)
	if err != nil {
		SetMsgInCache(
			data,
			qr.RequestID,
			porter.QFailed,
			"Internal error. Failed to construct north star.",
		)
		return
	}

	// Extra options below

	fieldsQ := qr.Fields
	regexQRaw := qr.Regex
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
				SetMsgInCache(
					data,
					qr.RequestID,
					porter.QFailed,
					"Failed to compile regexp.",
				)
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

	queryJob := QueryJob{
		ID:              qr.RequestID,
		IsSimpleQuery:   isSimpleQuery,
		Beacon:          beacon,
		ValidPrefix:     validPrefix,
		NorthStar:       northStar,
		FieldsQ:         fieldsQ,
		RegexQ:          regexQ,
		RegexQField:     regexQField,
		NegateUserRegex: negateUserRegex,
	}

	queueChan <- queryJob

	SetMsgInCache(
		data,
		qr.RequestID,
		porter.QWorking,
		"Query request passed validation.",
	)
}
