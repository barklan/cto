package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/caching"
	"github.com/barklan/cto/pkg/logserver"
	postgres "github.com/barklan/cto/pkg/postgres"
	"github.com/barklan/cto/pkg/restcore"
	"github.com/barklan/cto/pkg/storage"
	"github.com/jmoiron/sqlx"
)

func handleSysSignals(data *storage.Data) {
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
	data.InternalAlert(fmt.Sprintf("I received %s. Exiting now!", sigID))
	time.Sleep(200 * time.Millisecond)
	data.DB.Close()
	os.Exit(0)
}

func CrashExit(data *storage.Data, info string) {
	data.InternalAlert(fmt.Sprintf("Help! I crashed! %s", info))
	data.DB.Close()
	os.Exit(1)
}

func main() {
	log.Println("Starting...")

	// https://dgraph.io/docs/badger/faq/#are-there-any-go-specific-settings-that-i-should-use
	runtime.GOMAXPROCS(128)

	data := &storage.Data{}

	var rdb *sqlx.DB
	var err error
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

	db := storage.OpenDB("", "/main")
	data.DB = db
	defer db.Close()

	config := storage.ReadConfig(data)
	data.Config = config

	redis := caching.InitRedis()
	data.Cache = redis

	defer CrashExit(data, "Deferred in main.")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			CrashExit(data, "LogServer exited.")
			wg.Done()
		}()
		logserver.LogServerServe(data)
	}()

	wg.Add(1)
	go func() {
		defer func() {
			CrashExit(data, "chi server exited")
			wg.Done()
		}()
		restcore.Serve(data)
	}()

	// FIXME move to porter
	// tokenRotationTicker := time.NewTicker(4 * time.Hour)
	// go func() {
	// 	defer func() {
	// 		CrashExit(data, "Token rotation goroutine exited.")
	// 		wg.Done()
	// 	}()
	// 	for {
	// 		for projectName := range data.Config.P {
	// 			storage.RotateJWT(data, projectName)
	// 		}
	// 		<-tokenRotationTicker.C
	// 	}
	// }()

	go func() {
		handleSysSignals(data)
	}()

	data.InternalAlert("I am up!")

	wg.Wait()
	data.InternalAlert("All goroutines are done (or no one left alive). Main function will now exit.")
}
