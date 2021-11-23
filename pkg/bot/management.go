package bot

import (
	"fmt"
	"os"

	"github.com/barklan/cto/pkg/config"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func registerProjectManagementHandlers(b *tb.Bot, data *storage.Data) {
	b.Handle(tb.OnDocument, func(m *tb.Message) {
		if m.Sender.Username != data.SysAdmin {
			data.Send(data.Chat, fmt.Sprintf("@%s, help! Someone is trying to mess with me!", data.SysAdmin))
			return
		}

		fileName := m.Document.FileName
		if fileName[len(fileName)-4:] == ".yml" {
			configPath := "/app/config"
			if _, ok := os.LookupEnv("CTO_LOCAL_ENV"); ok {
				configPath = "environment"
			}

			b.Download(&m.Document.File, configPath+"/"+fileName)

			config := config.ReadConfig()
			data.Config = config
			data.CSend("Successfully reloaded configurations!")
		}
	})

	b.Handle("/projects", func(m *tb.Message) {
		if m.Sender.Username != data.SysAdmin {
			data.Send(data.Chat, fmt.Sprintf("@%s, help! Someone is trying to mess with me!", data.SysAdmin))
			return
		}

		data.CSend(fmt.Sprint(data.Config.EnvToProjectName))
	})
}