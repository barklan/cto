package bot

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

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
	secretKey := RandStringBytesMaskImpr(24) // TODO should be more secure and random
	prettyTitle := from.Title

	uid4, err := uuid.NewV4()
	if err != nil {
		s.logAndReport(from, "failed to generate uuid", err)
		return
	}
	u4 := uid4.String()

	s.Log.Info("opening tx to create new project", zap.String("projectTitle", prettyTitle))
	tx, err := s.R.Begin()
	if err != nil {
		s.logAndReport(from, "failed to create db transaction.", err)
		return
	}
	insert := "insert into project(id, client_id, pretty_title, secret_key) values ($1, $2, $3, $4)"
	if _, err = tx.Exec(insert, u4, client.ID, prettyTitle, secretKey); err != nil {
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
	configEnv := os.Getenv("CONFIG_ENV")
	coreHost := "cto_backend"
	if configEnv == "dev" {
		coreHost = "localhost"
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s:8888/api/core/setproject/%s", coreHost, u4),
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
		if m.Chat.Type == tb.ChatPrivate {
			s.JustSend(m.Chat, "User registration through telegrm not implemented.") // TODO
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

		email, ok, err := s.Cache.Get(m.Payload)
		if err != nil {
			s.JustSend(
				m.Chat,
				"Error when getting integration pass from cache.",
			)
			return
		}
		if !ok {
			s.JustSend(
				m.Chat,
				"Please visit link from ctopanel.com instead of directly invoking command.",
			)
			return
		}

		var existringEmail string
		err = s.R.Get(&existringEmail, "select email from client where tg_nick = $1", m.Sender.Username)
		if err != nil {
			s.Log.Error("failed to verify if tg nick already exists", zap.Error(err))
		}
		if existringEmail != "" && existringEmail != string(email) {
			s.JustSend(
				m.Chat,
				fmt.Sprintf(
					"This telegram user already registered for dirrerent email: <code>%s</code>.",
					existringEmail,
				),
				tb.ModeHTML,
			)
			return
		}

		if _, err := s.R.Exec(
			"update client set tg_nick = $1 where email = $2",
			m.Sender.Username,
			email,
		); err != nil {
			s.JustSend(m.Chat, "Failed to update client in database")
			return
		}

		client := models.Client{}
		if err := s.R.Get(&client, "select * from client where email = $1", email); err != nil {
			s.JustSend(
				m.Chat,
				"Could not get client with specified email.",
			)
			return
		}

		log.Println("Request to register new project.")
		s.newProject(m.Chat, &client)
	})
}
