package storage

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/dgraph-io/badger/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	Internal          = "Internal" // Reserved internal project - not listed in projects.
	internalKeySymbol = "!"
	variableKeySymbol = "$"
)

type Data struct {
	B         *tb.Bot
	Chat      *tb.Chat
	DB        *badger.DB
	Config    *Config
	MediaPath string
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

// TODO deprecate this
func (d *Data) GetStr(key string) string {
	varKey := variableKeySymbol + key
	return string(Get(d.DB, varKey))
}

func (d *Data) SetVar(projectName, key string, obj interface{}, ttl time.Duration) {
	var byteObj []byte
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.String {
		byteObj = []byte(obj.(string))
	} else {
		var err error
		byteObj, err = json.Marshal(obj)
		if err != nil {
			log.Panic(err)
		}
	}

	varKey := projectName + variableKeySymbol + key
	if ttl > 0 {
		SetWithTTL(d.DB, varKey, byteObj, ttl)
	} else {
		Set(d.DB, varKey, byteObj)
	}
}

func (d *Data) GetVar(projectName, key string) []byte {
	varKey := projectName + variableKeySymbol + key
	return Get(d.DB, varKey)
}

func (d *Data) VarExists(projectName, key string) bool {
	err := d.DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	switch err {
	case badger.ErrKeyNotFound:
		return false
	case nil:
		return true
	default:
		log.Panicln("panicing in VarExists function", err)
		return false
	}
}

func (d *Data) DeleteVar(projectName, key string) {
	err := d.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
	if err != nil {
		log.Panicln("panicing in DeleteVar function", err)
	}
}

// TODO deprecate this
func (d *Data) SetObj(key string, obj interface{}, ttl time.Duration) {
	var byteObj []byte
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.String {
		byteObj = []byte(obj.(string))
	} else {
		var err error
		byteObj, err = json.Marshal(obj)
		if err != nil {
			log.Panic(err)
		}
	}

	varKey := variableKeySymbol + key

	if ttl > 0 {
		SetWithTTL(d.DB, varKey, byteObj, ttl)
	} else {
		Set(d.DB, varKey, byteObj)
	}
}

// TODO deprecate this
func (d *Data) Get(key string) []byte {
	varKey := variableKeySymbol + key
	return Get(d.DB, varKey)
}

// Send sends a message and saves it to main storage.
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

// CSendSync sends to barklan with sync
func (d *Data) CSendSync(msg interface{}, options ...interface{}) (*tb.Message, error) {
	return d.Send(d.Chat, msg, options...)
}

// CSend sends to barklan
func (d *Data) CSend(msg interface{}, options ...interface{}) {
	go func() {
		_, _ = d.Send(d.Chat, msg, options...)
	}()
}

func (d *Data) PSend(projectName string, msg interface{}, options ...interface{}) {
	// TODO recovery for mute operation like in CSend or maybe fuck it?
	muted := d.VarExists(projectName, "muted")
	if muted {
		log.Println("I am muted!")
		return
	}

	chatID := d.Config.P[projectName]
	chat := &tb.Chat{ID: chatID}
	go func() {
		_, _ = d.Send(chat, msg, options...)
	}()
}
