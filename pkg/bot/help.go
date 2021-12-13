package bot

import (
	"fmt"

	"github.com/barklan/cto/pkg/storage"
	"github.com/barklan/cto/pkg/storage/vars"
	tb "gopkg.in/tucnak/telebot.v2"
)

func registerHelpHandler(data *storage.Data, b *tb.Bot) {
	b.Handle("/help", func(m *tb.Message) {
		// TODO should return project name, owner and secret.
		project, ok := VerifySender(data, m)
		if !ok {
			return
		}

		secret := data.GetVar(project, vars.SecretKey)
		owner := data.GetVar(project, vars.Owner)

		yourProjects := []string{}
		participatingIn := []string{}
		for p, cid := range data.Config.P {
			projectOwner := string(data.GetVar(p, vars.Owner))
			if projectOwner == m.Sender.Username {
				yourProjects = append(yourProjects, p)
			}

			cidChat, err := b.ChatByID(fmt.Sprint(cid))
			if err != nil {
				data.CSend(fmt.Sprintf("Something wrong with project %q", p))
			}

			_, err = b.ChatMemberOf(cidChat, m.Sender)
			if err == nil {
				participatingIn = append(participatingIn, p)
			}
		}

		// TODO enabled/disabled feature flags for projects
		data.PSend(project, fmt.Sprintf(
			`Project: <code>%s</code>; secret: <code>%s</code>; owner: <code>%s</code>.
Your projects: <code>%s</code>.
Participating in: <code>%s</code>.
Logging: on; alerting: on;`,
			project,
			secret,
			owner,
			yourProjects,
			participatingIn,
		), tb.ModeHTML)
	})
}
