package bot

import (
	"regexp"
	"strconv"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func registerOnTextHanler(b *tb.Bot, data *storage.Data) {
	b.Handle(tb.OnText, func(m *tb.Message) {
		if m.Chat.ID != data.Config.Internal.TG.BossChatID {
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
