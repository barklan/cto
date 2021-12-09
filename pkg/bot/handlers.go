package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/barklan/cto/pkg/storage"
	"github.com/barklan/cto/pkg/tempnginx"
	tb "gopkg.in/tucnak/telebot.v2"
)

var randSrc = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func registerTemporaryNginxHandler(b *tb.Bot, data *storage.Data) {
	b.Handle("/nginx", func(m *tb.Message) {
		projectName, ok := VerifySender(data, m)
		if !ok {
			return
		}

		minutes := 1
		basicAuthUsername := "nginx"
		basicAuthPassword := RandString(24)

		addressChan := make(chan string)
		go tempnginx.TemporaryNginx(
			data,
			projectName,
			addressChan,
			minutes,
			basicAuthUsername,
			basicAuthPassword,
		)

		address := <-addressChan

		wgetCmd := fmt.Sprintf(
			"`wget --no-verbose --no-parent --recursive --level=1 --no-directories --http-user=%s --http-password=%s %s`",
			basicAuthUsername,
			basicAuthPassword,
			address,
		)

		data.PSend(
			projectName,
			fmt.Sprintf(
				`Temporary nginx started at %s for %d minutes. To download all files use:
%s`,
				address,
				minutes,
				wgetCmd,
			),
			tb.ModeMarkdown,
		)
	})
}

func RegisterHandlers(b *tb.Bot, data *storage.Data) {
	registerStatusHandler(b, data)

	registerTemporaryNginxHandler(b, data)

	b.Handle("/start", func(m *tb.Message) {
		message := fmt.Sprintf(`ID of this chat: %s.
I will process requests only if this ID is set in configuration.
Your user ID is %s.
`,
			fmt.Sprint(m.Chat.ID), fmt.Sprint(m.Sender.ID))
		go func() {
			data.Send(m.Chat, message)
			b.Delete(m)
		}()
	})

	b.Handle("/clear", func(m *tb.Message) {
		if m.Sender.Username != data.SysAdmin {
			data.Send(data.Chat, fmt.Sprintf("@%s, help! Someone is trying to take my stuff!", data.SysAdmin))
			return
		}

		go func() {
			data.CSendSync("Started cleaning...")
			b.Delete(m)
			CleanUp(data)
		}()
	})

	b.Handle("/mute", func(m *tb.Message) {
		if m.Sender.Username != data.SysAdmin {
			data.Send(data.Chat, fmt.Sprintf("@%s, help! Someone is trying to shut me up!", data.SysAdmin))
			return
		}
		data.CSend("Muted for 4 hours.")
		data.SetObj("muted", "true", 4*time.Hour)
		b.Delete(m)
	})

	b.Handle("/unmute", func(m *tb.Message) {
		data.SetObj("muted", "false", 5*time.Minute)
		b.Delete(m)
		data.CSend("Unmuted.")
	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		// urls := []string{
		// "https://thatcopy.github.io/catAPI/imgs/jpg/60343c6.jpg",
		// "https://thatcopy.github.io/catAPI/imgs/jpg/9541262.jpg",
		// "https://thatcopy.github.io/catAPI/imgs/jpg/54adb30.jpg",
		// }

		AdminQuery(data, q)

		results := make(tb.Results, 1) // []tb.Result
		for i := 0; i < 1; i++ {
			result := &tb.ArticleResult{
				Title: q.Text,
				Text:  ".",
			}

			results[i] = result
			// needed to set a unique string ID for each result
			results[i].SetResultID(strconv.Itoa(i))

		}

		err := b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
		})
		if err != nil {
			log.Println(err)
		}
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

// Queries should only be allowed for responsible users.
func AdminQuery(data *storage.Data, q *tb.Query) error {
	if q.From.Username != data.SysAdmin {
		result := &tb.ArticleResult{
			Title: "You are not SysAdmin!",
			Text:  "",
		}
		result.SetResultID("1")
		err := data.B.Answer(q, &tb.QueryResponse{
			Results:   []tb.Result{result},
			CacheTime: 60, // a minute
		})
		if err != nil {
			log.Println(err)
		}

		return fmt.Errorf("Wrong query user.")
	}
	return nil
}
