package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/barklan/cto/pkg/bot"
	"github.com/barklan/cto/pkg/caching"
	porter "github.com/barklan/cto/pkg/porter"
	postgres "github.com/barklan/cto/pkg/postgres"
	"github.com/barklan/cto/pkg/storage"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
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
	log.Println(fmt.Sprintf("I received %s. Exiting now!", sigID))
	os.Exit(0)
}

// TODO handle closing database gracefully
func main() {
	log.Info("Starting...")

	config, err := storage.ReadInternalConfig("")
	if err != nil {
		log.Fatal(err)
	}

	var rdb *sqlx.DB
	for i := 0; i < 10; i++ {
		rdb, err = postgres.OpenDB()
		if err != nil {
			if i == 9 {
				log.Panicln("Could not open pg connection  10 times.")
			}
			time.Sleep(1 * time.Second)
			continue
		}
	}
	defer rdb.Close()

	base := porter.InitBase(&config, rdb)

	redis := caching.InitRedis()
	base.Cache = redis

	// TODO telebot migrating to v3 soon
	b := bot.Bot(config.TG.BotToken)

	sylon := bot.InitSylon(rdb, &config, b, redis)

	queries := make(chan porter.QueryRequestWrap, 10)
	go porter.Serve(base, sylon, queries)
	go porter.Publisher(base, queries)

	wg := new(sync.WaitGroup)

	tokenRotationTicker := time.NewTicker(4 * time.Hour)
	go func() {
		defer log.Panicln("Token rotation goroutine exited.")
		for {
			projects := make([]string, 0)
			if err := rdb.Select(&projects, "select id from project"); err != nil {
				log.Panicln("failed to get projects from db to rotate jwt")
			}

			for _, projectName := range projects {
				porter.RotateJWT(base, "guest", projectName)
			}
			<-tokenRotationTicker.C
		}
	}()

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
