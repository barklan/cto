package bot

import (
	"fmt"
	"log"

	"github.com/barklan/cto/pkg/postgres/models"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) registerHelpHandler() {
	s.B.Handle("/help", func(m *tb.Message) {
		// TODO should return project name, owner and secret.
		project, _, ok := s.VerifySender(m)
		if !ok {
			return
		}

		// FIXME select
		var owner models.Client
		if err := s.R.Get(&owner, "select * from client where id = $1", project.ClientID); err != nil {
			log.Panic("Owner must exist.")
		}

		// TODO recovery v5
		// yourProjects := []string{}
		// participatingIn := []string{}
		// for p, cid := range data.Config.P {
		// 	projectOwner := string(data.GetVar(p, vars.Owner))
		// 	if projectOwner == m.Sender.Username {
		// 		yourProjects = append(yourProjects, p)
		// 	}

		// 	cidChat, err := b.ChatByID(fmt.Sprint(cid))
		// 	if err != nil {
		// 		data.CSend(fmt.Sprintf("Something wrong with project %q", p))
		// 	}

		// 	_, err = b.ChatMemberOf(cidChat, m.Sender)
		// 	if err == nil {
		// 		participatingIn = append(participatingIn, p)
		// 	}
		// }

		// TODO enabled/disabled feature flags for projects
		s.PSend(project.ID, fmt.Sprintf(
			`Project: <code>%s</code>; secret: <code>%s</code>; owner: <code>%s</code>.`,
			project.ID,
			project.SecretKey,
			owner.TGNick,
		), tb.ModeHTML)
	})
}
