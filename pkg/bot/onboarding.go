package bot

import (
	"math/rand"
	"time"

	"github.com/barklan/cto/pkg/postgres/models"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) registerOnboardingHandlers() {
	s.B.Handle("/start", func(m *tb.Message) {
		if m.Sender.Username != "barklan" && m.Sender.Username != "qufiwefefwoyn" { // TODO
			s.JustSend(m.Chat, "You are not authorized to do anything.")
			return
		}

		if m.Chat.Type == tb.ChatPrivate {
			s.JustSend(m.Chat, "User registration not implemented.") // TODO
			return
		}

		if _, _, ok := s.VerifySender(m); ok {
			s.JustSend(
				m.Chat,
				"This chat is registered. Call <code>/help</code> for more info.",
				tb.ModeHTML,
			)
			return
		}

		client := models.Client{}
		if err := s.R.Get(&client, "select * from client where id = $1", m.Sender.ID); err != nil {
			s.JustSend(
				m.Chat,
				"This chat is not registered, but you cannot register projects "+
					"before you register yourself. To do that call <code>/start</code> "+
					"in personal chat with me.",
				tb.ModeHTML,
			)
			return
		}

		s.JustSend(m.Chat, "Project registration not implemented.")
		rand.Seed(time.Now().UnixNano())
		secretKey := RandStringBytesMaskImpr(48) // TODO should be more secure and random

		tx, err := s.R.Begin()
		if err != nil {
			s.JustSend(m.Chat, "Failed to create db transaction.")
			return
		}
		insert := "insert into project(client_id, secret_key) values ($1, $2)"
		if _, err = tx.Exec(insert, client.ID, secretKey); err != nil {
			s.JustSend(m.Chat, "Failed to insert new project.")
		}
	})
}
