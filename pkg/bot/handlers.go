package bot

import (
	"fmt"
	"time"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterHandlers(b *tb.Bot, data *storage.Data) {
	registerStatusHandler(b, data)

	b.Handle("/start", func(m *tb.Message) {
		message := fmt.Sprintf(`ID of this chat: %s.
I will process requests only if this ID is set in configuration.
Your user ID is %s.
`,
			fmt.Sprint(m.Chat.ID), fmt.Sprint(m.Sender.ID))
		go func() {
			data.Send(m.Chat, message)
		}()
	})

	// FIXME mutes everything - should be project specific
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

	registerOnTextHanler(b, data)

	registerProjectManagementHandlers(b, data)
}

// VerifySender returns projectName and if chat is registered
func VerifySender(data *storage.Data, m *tb.Message) (string, bool) {
	if v, ok := data.Config.ChatIDToProjectName[m.Chat.ID]; ok {
		return v, ok
	}
	data.JustSend(m.Chat, "I am not registered for this chat.")
	return "", false
}
