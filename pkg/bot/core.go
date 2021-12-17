package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/barklan/cto/pkg/porter"
	tb "gopkg.in/tucnak/telebot.v2"
)

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

func GetBoss(data *porter.Data) *tb.Chat {
	chatID := data.Config.TG.BossChatID
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
