package bot

import (
	"github.com/barklan/cto/pkg/postgres/models"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) RegisterHandlers() {
	// TODO recovery v5
	// registerStatusHandler(b, data)
	s.registerHelpHandler()
	// registerManagementHandlers(data, b)
	// registerOnTextHanler(data, b)

	// TODO recovery v5
	// b.Handle("/start", func(m *tb.Message) {
	// 	if m.Sender.Username != "barklan" { // TODO external auth service
	// 		data.JustSend(m.Chat, "You are not authorized to register projects.")
	// 		return
	// 	}

	// 	chatActual := m.Chat.ID
	// 	for project, chat := range data.Config.P {
	// 		if chatActual == chat {
	// 			data.JustSend(
	// 				m.Chat,
	// 				fmt.Sprintf("This chat is already registered for project %q", project),
	// 			)
	// 			return
	// 		}
	// 	}

	// 	rand.Seed(time.Now().UnixNano())

	// 	// TODO deprecated - uuid is used
	// 	// projectName := genUniqueProjectName(data)
	// 	secretKey := RandStringBytesMaskImpr(48) // TODO should be more secure and random

	// 	storage.AddProject(data, projectName, chatActual)
	// 	data.SetVar(projectName, vars.SecretKey, secretKey, -1)
	// 	data.SetVar(projectName, vars.Owner, m.Sender.Username, -1)

	// 	storage.RotateJWT(data, projectName)

	// 	data.JustSend(
	// 		m.Chat,
	// 		fmt.Sprintf(
	// 			"Success! Project %q registered for this chat. Your secret key:\n `%s`",
	// 			projectName,
	// 			secretKey,
	// 		),
	// 		tb.ModeMarkdown,
	// 	)
	// })

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
		s.JustSend(m.Chat, "This chat is not registered for any project.")
		return nil, nil, false
	}

	project := &models.Project{}
	s.R.Get(project, "select * from project where id = $1", chat.ProjectID)

	return project, chat, true
}
