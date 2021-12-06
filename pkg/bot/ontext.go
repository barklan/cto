package bot

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func registerOnTextHanler(b *tb.Bot, data *storage.Data) {
	b.Handle(tb.OnText, func(m *tb.Message) {
		// TODO boss only
		if m.Sender.Username != data.SysAdmin {
			data.Send(data.Chat, fmt.Sprintf("@%s, help! Someone is trying to mess with me!", data.SysAdmin))
			return
		}

		r, _ := regexp.Compile(`^Set error similarity threshold (0\.\d+)$`)
		arr := r.FindStringSubmatch(m.Text)
		if arr != nil {
			simStr := arr[1]
			simFloat, err := strconv.ParseFloat(simStr, 64)
			if err != nil {
				data.CSend("Failed to parse float.")
			}

			data.Config.Internal.Log.SimilarityThreshold = simFloat
			data.CSend("New threshold is set.")
		}
	})
}
