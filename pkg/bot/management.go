package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) bossChat(m *tb.Message) bool {
	return m.Chat.ID == s.Config.TG.BossChatID
}
