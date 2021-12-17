package porter

import (
	"fmt"
	"log"
	"os"

	"github.com/barklan/cto/pkg/storage"
	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Data struct {
	B         *tb.Bot
	Chat      *tb.Chat
	Config    *storage.InternalConfig
	MediaPath string
	R         *sqlx.DB
}

func InitData() *Data {
	data := Data{}

	configEnvironment, ok := os.LookupEnv("CONFIG_ENV")
	if !ok {
		log.Panic("Config environment variable CONFIG_ENV must be specified.")
	}
	if configEnvironment == "dev" {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}

		data.MediaPath = currentDir + "/.cache/media"
	} else {
		data.MediaPath = "/app/media"
	}

	return &data
}


// CreateMediaDirIfNotExists creates the directory in default media path.
// It can accept nested directory path, but all parent directories must
// exist. Returns full directory path.
func (d *Data) CreateMediaDirIfNotExists(dirname string) string {
	fullDirname := d.MediaPath + "/" + dirname
	_, err := os.Stat(fullDirname)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(fullDirname, 0755)
		if errDir != nil {
			log.Panic(err)
		}
	}

	return fullDirname
}

func (d *Data) Send(to tb.Recipient, msg interface{}, options ...interface{}) (*tb.Message, error) {
	m, err := d.B.Send(to, msg, options...)
	if err != nil {
		log.Printf("Failed to send tg message. %v", err)
		return nil, err
	}
	log.Printf("Send TG message %v\n", msg)
	return m, err
}

func (d *Data) JustSend(to tb.Recipient, msg interface{}, options ...interface{}) {
	go func() {
		_, _ = d.Send(to, msg, options...)
	}()
}

func (d *Data) CSendSync(msg interface{}, options ...interface{}) (*tb.Message, error) {
	return d.Send(d.Chat, msg, options...)
}

func (d *Data) CSend(msg interface{}, options ...interface{}) {
	go func() {
		_, _ = d.Send(d.Chat, msg, options...)
	}()
}

func (d *Data) PSend(projectName string, msg interface{}, options ...interface{}) {
	// TODO recovery for mute operation for v5
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
