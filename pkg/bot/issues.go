package bot

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) NotifyAboutError(
	projectID, env, service, timestamp, key, flag string,
) {
	queryString := url.QueryEscape(key)
	exactLogURL := fmt.Sprintf(
		"%s/log/exact?key=%s",
		s.Config.Log.ServiceHostname,
		queryString,
	)
	selector := &tb.ReplyMarkup{}
	btnURL := selector.URL("View Log", exactLogURL)
	selector.Inline(
		selector.Row(btnURL),
	)

	var extraTimeStr string
	fluentdTime, err := time.Parse("2006-01-02 15:04:05 -0700", timestamp)
	if err == nil {
		now := time.Now()
		elapsed := now.Sub(fluentdTime).Round(time.Second)
		extraTimeStr = fmt.Sprintf(" (%s ago)", elapsed)
	}

	upperFlag := strings.ToUpper(flag)
	message := fmt.Sprintf(
		"*%s* in %s (%s) at %s%s.",
		upperFlag,
		service,
		env,
		timestamp,
		extraTimeStr,
	)

	s.PSend(projectID, message, tb.ModeMarkdown, selector)
}
