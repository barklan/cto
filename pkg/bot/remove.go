package bot

import (
	"database/sql"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) registerRemoveHandler() {
	s.B.Handle("/remove", func(m *tb.Message) {
		project, _, ok := s.VerifySender(m)
		if !ok {
			return
		}

		var ownerUsername sql.NullString
		if err := s.R.Get(
			&ownerUsername,
			`--sql
				select tg_nick
				from client inner join project
					on project.client_id = client.id
				where project.id = $1`,
			project.ID,
		); err != nil {
			s.JustSend(m.Chat, "Failed to verify owner of the project.")
			return
		}

		if !ownerUsername.Valid {
			s.JustSend(m.Chat, "Failed to verify owner of the project.")
			return
		}

		if m.Sender.Username != ownerUsername.String {
			s.JustSend(m.Chat, "You are not the owner of this project.")
			return
		}

		if _, err := s.R.Exec("delete from project where id = $1", project.ID); err != nil {
			s.JustSend(m.Chat, "Failed to remove project.")
			return
		}
		s.JustSend(m.Chat, "Project removed.")
	})
}
