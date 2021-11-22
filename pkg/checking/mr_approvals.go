package checking

import (
	"fmt"
	"log"
	"time"

	"github.com/barklan/cto/pkg/gitlab"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func getEventMrIidsInWindow(after, before string) ([]int64, error) {
	var query string
	if before == "" {
		query = fmt.Sprintf("&after=%s", after)
	} else {
		query = fmt.Sprintf("&after=%s&before=%s", after, before)
	}

	events, err := gitlab.GetApprovedMREvents(query)
	if err != nil {
		return nil, err
	}
	eventMrIids := make([]int64, 0)
	for _, event := range events {
		eventMrIids = append(eventMrIids, event.TargetIid)
	}
	return eventMrIids, nil
}

func CheckEachMrInApprovedMREventsSlow(b *tb.Bot, data *storage.Data, args ...interface{}) {
	now := time.Now()
	layout := "2006-01-02"

	aggregateIids := make([]int64, 0)
	howManyWeeksToCheck := 10
	for i := 0; i < howManyWeeksToCheck; i++ {
		days := i * 7
		cursorTime := now.Add(time.Duration(-days*24) * time.Hour)
		cursorTimeStr := cursorTime.Format(layout)

		aWeekAgo := now.Add(time.Duration(-(days+7)*24) * time.Hour)
		aWeekAgoStr := aWeekAgo.Format(layout)

		events, err := getEventMrIidsInWindow(aWeekAgoStr, cursorTimeStr)
		if err != nil {
			data.CSend("Failed to get gitlab events.")
		}
		aggregateIids = append(aggregateIids, events...)
	}

	eventMrIids := unique(aggregateIids)

	openMrs, err := gitlab.GetOpenMergeRequests("")
	if err != nil {
		data.CSend("Failed to get open merge requests.")
	}
	openMrIids := make([]int64, 0)
	for _, mr := range openMrs {
		openMrIids = append(openMrIids, mr.Iid)
	}

	intersection := Intersection(eventMrIids, openMrIids)
	if len(intersection) > 0 {
		data.CSend(fmt.Sprintf("The following merge requests have "+
			"been approved in past %d weeks but still open: %s. "+
			"This makes me sad. :(", howManyWeeksToCheck, fmt.Sprint(intersection)))
	}
}

func CheckEachMrInApprovedMREventsFast(b *tb.Bot, data *storage.Data, args ...interface{}) {
	layout := "2006-01-02"

	now := time.Now()

	lastDay := now.Add(-24 * time.Hour)
	lastDayStr := lastDay.Format(layout)

	events, err := getEventMrIidsInWindow(lastDayStr, "")
	if err != nil {
		data.CSend("Failed to get gitlab events.")
	}
	log.Println(events)

	openMrs, err := gitlab.GetOpenMergeRequests("")
	if err != nil {
		data.CSend("Failed to get open merge requests.")
	}
	openMrIids := make([]int64, 0)
	for _, mr := range openMrs {
		openMrIids = append(openMrIids, mr.Iid)
	}

	intersection := Intersection(events, openMrIids)
	if len(intersection) > 0 && data.GetStr("fastMRChecksMuted") == "false" {
		selector := &tb.ReplyMarkup{ReplyKeyboardRemove: true}
		btnMute := selector.Data("Mute", "mute")
		selector.Inline(
			selector.Row(btnMute),
		)

		data.CSend(
			fmt.Sprintf("The following merge requests have "+
				"been approved today but still open: %s.", fmt.Sprint(intersection)),
			selector,
		)

		b.Handle(&btnMute, func(c *tb.Callback) {
			b.Respond(c, &tb.CallbackResponse{Text: "Muted."})

			muteFor := 6 * time.Hour
			data.SetObj("fastMRChecksMuted", "true", muteFor)

			go func() {
				time.Sleep(muteFor)
				data.SetObj("fastMRChecksMuted", "false", -1)
			}()
		})

	}
}

func Intersection(a, b []int64) (c []int64) {
	m := make(map[int64]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func unique(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func mergeSlices(args ...[]int64) []int64 {
	mergedSlice := make([]int64, 0)
	for _, oneSlice := range args {
		mergedSlice = append(mergedSlice, oneSlice...)
	}

	return mergedSlice
}
