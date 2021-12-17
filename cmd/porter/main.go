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

func handleSysSignals(base *porter.Base, sylon *bot.Sylon) {
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
	sylon.B.Close()
	base.R.Close()
	log.Println(fmt.Sprintf("I received %s. Exiting now!", sigID))
	os.Exit(0)
}

func main() {
	log.Println("Starting...")

	config, err := storage.ReadInternalConfig("")
	if err != nil {
		log.Fatal(err)
	}

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

	base := porter.InitBase(&config, rdb)

	// TODO telebot migrating to v3 soon
	b := bot.Bot(config.TG.BotToken)

	sylon := bot.InitSylon(rdb, &config, b)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		base.ServeGRPC(sylon)
	}()

	wg.Add(1)
	go func() {
		defer func() {
			log.Panic("Telebot poller exited.")
			wg.Done()
		}()

		sylon.RegisterHandlers()
		b.Start()
	}()

	go func() {
		handleSysSignals(base, sylon)
	}()

	wg.Wait()
}
