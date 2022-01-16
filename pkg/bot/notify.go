package bot

import (
	"fmt"
	"log"
	"os"
)

type Chat struct{}

// DEPRECATED. This relies on env variables.
func Notify(message string) {
	chatID := os.Getenv("TELEGRAM_GROUP_CHAT_ID")
	_, err := request(fmt.Sprintf("sendMessage?chat_id=%s&text=%s", chatID, message))
	if err != nil {
		log.Printf("Failed to send message")
	}
}
