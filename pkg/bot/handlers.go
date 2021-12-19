package bot

import (
	"github.com/barklan/cto/pkg/postgres/models"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) RegisterHandlers() {
	// TODO recovery v5
	// registerStatusHandler(b, data)
	s.registerOnboardingHandlers() // `/start`
	s.registerHelpHandler()        // `/help`
	s.registerRemoveHandler()      // `/remove`
	// registerManagementHandlers(data, b)
	// registerOnTextHanler(data, b)

	// TODO recoery v5
	// b.Handle("/mute", func(m *tb.Message) {
	// 	project, ok := VerifySender(data, m)
	// 	if !ok {
	// 		return
	// 	}
	// 	data.PSend(project, "Muted for 4 hours.")
	// 	data.SetVar(project, "muted", "flag", 4*time.Hour)
	// })

	// b.Handle("/unmute", func(m *tb.Message) {
	// 	project, ok := VerifySender(data, m)
	// 	if !ok {
	// 		return
	// 	}
	// 	data.DeleteVar(project, "muted")
	// 	data.CSend("Unmuted.")
	// })

	// TODO
	// registerOnTextHanler(b, data)
}

// VerifySender returns projectName and if chat is registered
func (s *Sylon) VerifySender(m *tb.Message) (*models.Project, *models.Chat, bool) {
	chat := &models.Chat{}
	if err := s.R.Get(chat, "select * from chat where id = $1", m.Chat.ID); err != nil {
		return nil, nil, false
	}

	project := &models.Project{}
	if err := s.R.Get(project, "select * from project where id = $1", chat.ProjectID); err != nil {
		s.CSend("Fail in VerifySender.")
		return nil, nil, false
	}

	return project, chat, true
}
