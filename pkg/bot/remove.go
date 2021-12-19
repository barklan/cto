package bot

import tb "gopkg.in/tucnak/telebot.v2"

func (s *Sylon) registerRemoveHandler() {
	s.B.Handle("/remove", func(m *tb.Message) {
		project, chat, ok := s.VerifySender(m)
		if !ok {
			return
		}

		s.R.MustExec("delete from project ", args ...interface{})
	})
}
