package bot

import (
	"fmt"
	"log"

	"github.com/barklan/cto/pkg/postgres/models"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) registerHelpHandler() {
	s.B.Handle("/help", func(m *tb.Message) {
		s.JustSend(m.Chat, fmt.Sprintf(
			"Chat ID: <code>%d</code>; Your ID: <code>%d</code>.",
			m.Chat.ID, m.Sender.ID,
		), tb.ModeHTML)

		// TODO should return project name, owner and secret.
		if m.Chat.Type == tb.ChatPrivate {
			var client models.Client
			if err := s.R.Get(&client, "select * from client where personal_chat = $1", m.Chat.ID); err != nil {
				s.JustSend(m.Chat,
					"Personal chat. You are not registered. "+
						"To register call <code>/start</code>.",
					tb.ModeHTML,
				)
				return
			}

			s.JustSend(m.Chat, "Personal chat. You are registered.")
			return
		}

		project, _, ok := s.VerifySender(m)
		if !ok {
			s.JustSend(
				m.Chat,
				"Group chat. No project is registered for this chat. "+
					"To start a new project for this group call <code>/start</code>.",
				tb.ModeHTML,
			)
			return
		}

		// FIXME select
		var owner models.Client
		if err := s.R.Get(&owner, "select * from client where id = $1", project.ClientID); err != nil {
			log.Panic("Owner must exist.")
		}

		// TODO this should be visible in personal chats
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
			`Group chat. Registered project: <code>%s</code>; secret: <code>%s</code>; owner: <code>%s</code>.`,
			project.ID,
			project.SecretKey,
			owner.TGNick,
		), tb.ModeHTML)
	})
}
