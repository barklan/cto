package bot

import tb "gopkg.in/tucnak/telebot.v2"

func (s *Sylon) registerRemoveHandler() {
	s.B.Handle("/remove", func(m *tb.Message) {
		project, _, ok := s.VerifySender(m)
		if !ok {
			return
		}

		var ownerChatID int64
		if err := s.R.Get(
			&ownerChatID,
			`--sql
			select personal_chat
			from client inner join project
				on project.client_id = client.id
			where project.id = $1`,
			project.ID,
		); err != nil {
			s.JustSend(m.Chat, "Failed to verify owner of the project.")
		}

		if m.Sender.ID != ownerChatID {
			s.JustSend(m.Chat, "You are not the owner of this project.")
		}

		if _, err := s.R.Exec("delete from project where id = $1", project.ID); err != nil {
			s.JustSend(m.Chat, "Failed to remove project.")
		}
		s.JustSend(m.Chat, "Project removed.")
	})
}
