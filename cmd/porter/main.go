package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/barklan/cto/pkg/bot"
	porter "github.com/barklan/cto/pkg/porter"
	postgres "github.com/barklan/cto/pkg/postgres"
	"github.com/barklan/cto/pkg/storage"
	"github.com/jmoiron/sqlx"
)

func handleSysSignals(data *porter.Data) {
	SigChan := make(chan os.Signal, 1)

	signal.Notify(SigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-SigChan
	var sigID string
	switch sig {
	case syscall.SIGHUP:
		sigID = "SIGHUP"
	case syscall.SIGINT:
		sigID = "SIGINT"
	case syscall.SIGTERM:
		sigID = "SIGTERM"
	case syscall.SIGQUIT:
		sigID = "SIGQUIT"
	default:
		sigID = "UNKNOWN"
	}
	data.B.Close()
	log.Println(fmt.Sprintf("I received %s. Exiting now!", sigID))
	os.Exit(0)
}

func main() {
	log.Println("Starting...")

	data := porter.InitData()
	config, err := storage.ReadInternalConfig("")
	if err != nil {
		log.Fatal(err)
	}
	data.Config = &config

	// TODO telebot migrating to v3 soon
	b := bot.Bot(config.TG.BotToken)
	data.B = b

	data.Chat = bot.GetBoss(data)

	var rdb *sqlx.DB
	for i := 0; i < 10; i++ {
		rdb, err = postgres.OpenDB()
		if err != nil {
			if i == 9 {
				log.Panicf("Could not open pg connection  10 times.")
			}
			time.Sleep(1 * time.Second)
			continue
		}

	}

	data.R = rdb

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		porter.Serve(data)
	}()

	wg.Add(1)
	go func() {
		defer func() {
			log.Panic("Telebot poller exited.")
			wg.Done()
		}()

		bot.RegisterHandlers(b, data)
		b.Start()
	}()

	go func() {
		handleSysSignals(data)
	}()

	wg.Wait()
}
