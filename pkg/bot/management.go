package bot

import (
	"github.com/barklan/cto/pkg/porter"
	tb "gopkg.in/tucnak/telebot.v2"
)

func bossChat(data *porter.Data, m *tb.Message) bool {
	return m.Chat.ID == data.Config.TG.BossChatID
}
