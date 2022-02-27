package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/caching"
	"github.com/barklan/cto/pkg/core/storage"
	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Sylon struct {
	Log    *zap.Logger
	R      *sqlx.DB
	Config *storage.InternalConfig
	B      *tb.Bot
	Chat   *tb.Chat
	Cache  caching.Cache
}

func InitSylon(
	r *sqlx.DB,
	config *storage.InternalConfig,
	b *tb.Bot,
	cache caching.Cache,
	lg *zap.Logger,
) *Sylon {
	chatID := config.TG.BossChatID
	chat, err := b.ChatByID(fmt.Sprint(chatID))
	if err != nil {
		log.Panicln("failed to init sylon", err)
	}
	sylon := &Sylon{
		Log:    lg,
		R:      r,
		Config: config,
		B:      b,
		Chat:   chat,
		Cache:  cache,
	}
	return sylon
}

func (s *Sylon) Send(
	to tb.Recipient,
	msg string,
	options ...interface{},
) (*tb.Message, error) {
	m, err := s.B.Send(to, msg, options...)
	if err != nil {
		s.Log.Error("failed to send tg message", zap.String("msg", msg), zap.Error(err))
		return nil, err
	}
	s.Log.Info("sent tg msg", zap.String("msg", msg))
	return m, err
}

func (d *Sylon) JustSend(to tb.Recipient, msg string, options ...interface{}) {
	go func() {
		_, _ = d.Send(to, msg, options...)
	}()
}

func (d *Sylon) CSend(msg string, options ...interface{}) {
	go func() {
		_, _ = d.Send(d.Chat, msg, options...)
	}()
}

func (d *Sylon) PSend(projectName string, msg string, options ...interface{}) {
	// TODO recovery for mute operation for v5 and use this instead of JustSend where it is meant to be
	// muted := d.VarExists(projectName, "muted")
	// if muted {
	// 	log.Println("I am muted!")
	// 	return
	// }

	var chatID int64
	if err := d.R.Get(&chatID, "SELECT id FROM chat WHERE project_id=$1", projectName); err != nil {
		d.CSend(fmt.Sprintf("chat_id not found for project %q", projectName))
	}
	chat := &tb.Chat{ID: chatID}
	go func() {
		_, _ = d.Send(chat, msg, options...)
	}()
}

func Bot(botToken string) *tb.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panic(err)
		return nil
	}
	return b
}

func GetBoss(sylon *Sylon) *tb.Chat {
	chatID := sylon.Config.TG.BossChatID
	return &tb.Chat{ID: chatID}
}

// Direct request to telegram api. Should not be used.
func request(botMethod string) ([]byte, error) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	req, _ := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://api.telegram.org/bot%s/%s", botToken, botMethod), nil)
	log.Println("Sending request to telegram api.")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Telegram client get failed: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	log.Println(string(body))
	return body, nil
}
