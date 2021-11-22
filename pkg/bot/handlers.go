package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/barklan/cto/pkg/gitlab"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func RegisterHandlers(b *tb.Bot, data *storage.Data) {
	registerStatusHandler(b, data)

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

	b.Handle("/help", func(m *tb.Message) {
		menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		readmeBtn := menu.URL("Project Readme", gitlab.GetProjectReadme())
		menu.Inline(
			menu.Row(readmeBtn),
		)
		b.Delete(m)
		data.CSend(
			getHelpString(data), menu, tb.ModeMarkdown)
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
}

func getRequest(url string) string {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GITLAB_API_TOKEN")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

// VerifySender returns projectName and if chat is registered
func VerifySender(data *storage.Data, m *tb.Message) (string, bool) {
	if v, ok := data.Config.ChatIDToProjectName[m.Chat.ID]; ok {
		return v, ok
	}
	data.JustSend(m.Chat, "I am not registered for this chat.")
	return "", false
}

// TODO deprecate
// func OnlyMainChat(data *storage.Data, m *tb.Message) error {
// if m.Chat.ID != data.Chat.ID {
// data.Send(m.Chat, "I am not registered for this chat.")
// return fmt.Errorf("Wrong chat")
// }
// return nil
// }

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
