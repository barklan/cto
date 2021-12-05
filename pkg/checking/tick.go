package checking

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Blocking function.
func Tick(b *tb.Bot, data *storage.Data, duration time.Duration, checkerFunc func(*tb.Bot, *storage.Data, ...interface{}), args ...interface{}) {
	ticker := time.NewTicker(duration)
	log.Print("Ticker created.")

	for {
		<-ticker.C
		checkerFunc(b, data, args...)
	}
}

func GoCheck(
	b *tb.Bot,
	data *storage.Data,
	wg *sync.WaitGroup,
	checkTitle string,
	interval time.Duration,
	checkerFunc func(*tb.Bot, *storage.Data, ...interface{}),
	args ...interface{},
) {
	wg.Add(1)
	go func() {
		defer func() {
			data.Send(data.Chat, fmt.Sprintf("%s crashed. @%s.", checkTitle, data.SysAdmin))
			wg.Done()
		}()

		log.Printf("GoCheck invoked for %q.", checkTitle)

		Tick(b, data, interval, checkerFunc, args...)
	}()
}
