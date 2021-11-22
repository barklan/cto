package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

func GetBoss(chatID int64) *tb.Chat {
	return &tb.Chat{ID: chatID}
}

// TODO recovery
// func FindMainChatCreatorUsername(data *storage.Data) string {
// 	admins, err := data.B.AdminsOf(data.Chat)
// 	if err != nil {
// 		return "somebody"
// 	}

// 	for _, admin := range admins {
// 		if admin.Role == tb.Creator {
// 			return admin.User.Username
// 		}
// 	}
// 	return "somebody"
// }

// func BegForAdminRights(data *storage.Data) {
// 	log.Println("Checking for admin rights...")
// 	admins, err := data.B.AdminsOf(data.Chat)
// 	if err != nil {
// 		return
// 	}

// 	_ = FindMainChatCreatorUsername(data)

// 	for _, admin := range admins {
// 		if admin.User.ID == data.B.Me.ID {
// 			if admin.CanDeleteMessages == false {
// 				data.CSend(fmt.Sprintf(
// 					"Please set me permission to delete messages " +
// 						"to clean up after myself.",
// 					// chatCreator,
// 				))
// 			}
// 			log.Println("I am admin!")
// 			return
// 		}
// 	}

// 	data.CSend(fmt.Sprintf(
// 		"Please make me admin of this chat to clean up after myself.\n" +
// 			"If you do so, please change chat id in configuration after that.",
// 		// chatCreator,
// 	))
// 	log.Println("Begging.. for admin")
// }

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
