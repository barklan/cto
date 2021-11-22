package checking

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func CheckByExternalRequest(b *tb.Bot, data *storage.Data, args ...interface{}) {
	url := args[0].(string)
	try := args[1].(int)
	badgerKey := fmt.Sprintf("isHalted-%s", url)
	isHaltedStr := data.GetStr(badgerKey)

	isHalted := false
	if isHaltedStr != "" {
		var err error
		isHalted, err = strconv.ParseBool(isHaltedStr)
		if err != nil {
			log.Println("failed to parse boolish value:", err)
		}
	}

	log.Printf("checking %s, try %d", url, try)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if try == 0 && isHalted == false {
			time.Sleep(30 * time.Second)
			CheckByExternalRequest(b, data, url, try+1)
			return
		}
		if isHalted == false {
			data.CSend(fmt.Sprintf("Cannot reach:\n%s.", url))
			data.SetObj(badgerKey, "true", 45*time.Minute)
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if try == 0 && isHalted == false {
			time.Sleep(30 * time.Second)
			CheckByExternalRequest(b, data, url, try+1)
			return
		}
		if isHalted == false {
			msg := fmt.Sprintf("Code %d:\n%s", resp.StatusCode, url)
			data.CSend(msg)
			data.SetObj(badgerKey, "true", 45*time.Minute)
		}
		return
	}
	if isHalted == true {
		data.SetObj(badgerKey, "false", 5*time.Minute)
		data.CSend(fmt.Sprintf("%s is up again!", url))
	}
}
