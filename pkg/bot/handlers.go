package bot

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/barklan/cto/pkg/storage"
	"github.com/barklan/cto/pkg/storage/vars"
	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterHandlers(b *tb.Bot, data *storage.Data) {
	registerStatusHandler(b, data)
	registerHelpHandler(data, b)
	registerManagementHandlers(data, b)
	registerOnTextHanler(data, b)

	b.Handle("/start", func(m *tb.Message) {
		if m.Sender.Username != "barklan" { // TODO external auth service
			data.JustSend(m.Chat, "You are not authorized to register projects.")
			return
		}

		chatActual := m.Chat.ID
		for project, chat := range data.Config.P {
			if chatActual == chat {
				data.JustSend(
					m.Chat,
					fmt.Sprintf("This chat is already registered for project %q", project),
				)
				return
			}
		}

		rand.Seed(time.Now().UnixNano())

		projectName := genUniqueProjectName(data)
		secretKey := RandStringBytesMaskImpr(48) // TODO should be more secure and random

		storage.AddProject(data, projectName, chatActual)
		data.SetVar(projectName, vars.SecretKey, secretKey, -1)
		data.SetVar(projectName, vars.Owner, m.Sender.Username, -1)

		storage.RotateJWT(data, projectName)

		data.JustSend(
			m.Chat,
			fmt.Sprintf(
				"Success! Project %q registered for this chat. Your secret key:\n `%s`",
				projectName,
				secretKey,
			),
			tb.ModeMarkdown,
		)
	})

	b.Handle("/mute", func(m *tb.Message) {
		project, ok := VerifySender(data, m)
		if !ok {
			return
		}
		data.PSend(project, "Muted for 4 hours.")
		data.SetVar(project, "muted", "flag", 4*time.Hour)
	})

	b.Handle("/unmute", func(m *tb.Message) {
		project, ok := VerifySender(data, m)
		if !ok {
			return
		}
		data.DeleteVar(project, "muted")
		data.CSend("Unmuted.")
	})

	// TODO
	// registerOnTextHanler(b, data)
}

// VerifySender returns projectName and if chat is registered
func VerifySender(data *storage.Data, m *tb.Message) (string, bool) {
	for project, chatID := range data.Config.P {
		if chatID == m.Chat.ID {
			return project, true
		}
	}
	data.JustSend(m.Chat, "No project is registered for this chat.")
	return "", false
}
