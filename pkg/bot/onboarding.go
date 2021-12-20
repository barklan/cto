package bot

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/postgres/models"
	"github.com/gofrs/uuid"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (s *Sylon) logAndReport(chat *tb.Chat, msg string, err error) {
	e := fmt.Errorf("%s: %w", msg, err)
	log.Println(e)
	s.JustSend(chat, msg)
}

func (s *Sylon) newProject(from *tb.Chat, client *models.Client) {
	rand.Seed(time.Now().UnixNano())
	secretKey := RandStringBytesMaskImpr(48) // TODO should be more secure and random

	uid4, err := uuid.NewV4()
	if err != nil {
		s.logAndReport(from, "failed to generate uuid", err)
		return
	}
	u4 := uid4.String()

	log.Println("opening tx to create new project")
	tx, err := s.R.Begin()
	if err != nil {
		s.logAndReport(from, "failed to create db transaction.", err)
		return
	}
	insert := "insert into project(id, client_id, secret_key) values ($1, $2, $3)"
	if _, err = tx.Exec(insert, u4, client.ID, secretKey); err != nil {
		s.logAndReport(from, "failed to insert new project", err)
		if e := tx.Rollback(); e != nil {
			s.logAndReport(from, "failed to rollback transaction", err)
		}
		return
	}

	insert = "insert into chat(id, project_id) values ($1, $2)"
	if _, err = tx.Exec(insert, from.ID, u4); err != nil {
		s.logAndReport(from, "failed to insert chat", err)
		if e := tx.Rollback(); e != nil {
			s.logAndReport(from, "failed to rollback transaction", err)
		}
		return
	}

	log.Println("requesting cto-core to acknoledge project")
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://cto_backend:8888/api/core/setproject/%s", u4),
		nil,
	)
	if err != nil {
		s.logAndReport(from, "failed to create request", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logAndReport(from, "failed to propagate project to one of core replicas", err)
		if e := tx.Rollback(); e != nil {
			s.logAndReport(from, "failed to rollback", err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logAndReport(
			from,
			fmt.Sprintf(
				"core replica denied new project request (status code %d)",
				resp.StatusCode,
			),
			err,
		)
		if e := tx.Rollback(); e != nil {
			s.logAndReport(from, "failed to rollback", err)
		}
		return
	}

	if err := tx.Commit(); err != nil {
		s.logAndReport(from, "failed to commit tx", err)
	}

	s.JustSend(from, fmt.Sprintf(
		"New project <code>%s</code> created with secret <code>%s</code>.",
		u4, secretKey,
	), tb.ModeHTML)
}

func (s *Sylon) registerOnboardingHandlers() {
	s.B.Handle("/start", func(m *tb.Message) {
		if m.Sender.Username != "barklan" &&
			m.Sender.Username != "qufiwefefwoyn" { // TODO
			s.JustSend(m.Chat, "You are not authorized to do anything.")
			return
		}

		if m.Chat.Type == tb.ChatPrivate {
			s.JustSend(m.Chat, "User registration not implemented.") // TODO
			return
		}

		if _, _, ok := s.VerifySender(m); ok {
			s.JustSend(
				m.Chat,
				"This chat is registered. Call <code>/help</code> for more info.",
				tb.ModeHTML,
			)
			return
		}

		client := models.Client{}
		if err := s.R.Get(&client, "select * from client where personal_chat = $1", m.Sender.ID); err != nil {
			s.JustSend(
				m.Chat,
				"This chat is not registered, but you cannot register projects "+
					"before you register yourself. To do that call <code>/start</code> "+
					"in personal chat with me.",
				tb.ModeHTML,
			)
			return
		}

		log.Println("Request to register new project.")
		s.newProject(m.Chat, &client)
	})
}
