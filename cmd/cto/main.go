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
	"github.com/barklan/cto/pkg/checking"
	"github.com/barklan/cto/pkg/config"
	"github.com/barklan/cto/pkg/grpcsrv"
	"github.com/barklan/cto/pkg/logserver"
	"github.com/barklan/cto/pkg/storage"
	"github.com/golang-jwt/jwt/v4"
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
	data.CSendSync(fmt.Sprintf("I received %s. Exiting now!", sigID))
	time.Sleep(200 * time.Millisecond)
	data.DB.Close()
	data.B.Close()
	os.Exit(0)
}

func CrashExit(data *storage.Data, info string) {
	data.CSendSync(fmt.Sprintf("Help! I crashed! %s", info))
	data.DB.Close()
	data.B.Close()
	os.Exit(1)
}

func main() {
	log.Println("Starting...")

	config := config.ReadConfig()

	data := storage.InitData()

	data.Config = config

	db := storage.OpenDB("", "/main")
	data.DB = db
	defer db.Close()

	b := bot.Bot(config.Internal.TG.BotToken)
	data.B = b

	data.Chat = bot.GetBoss(data)

	// TODO get rid of this after getting rid of grpc server
	storage.GData = data

	defer CrashExit(data, "Deferred in main.")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			CrashExit(data, "Telebot poller exited.")
			wg.Done()
		}()

		bot.RegisterHandlers(b, data)
		b.Start()
	}()

	wg.Add(1)
	go func() {
		defer func() {
			CrashExit(data, "All checks exited.")
			wg.Done()
		}()
		for projectName := range data.Config.P {
			time.Sleep(2 * time.Second) // we want some interval between outgoing requests
			checking.LaunchChecks(b, data, projectName)
		}
	}()

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
			data.CSend("grpc server deferred called")
			wg.Done()
		}()

		grpcsrv.Serve(data)
	}()

	tokenRotationTicker := time.NewTicker(4 * time.Hour)
	go func() {
		defer func() {
			CrashExit(data, "Token rotation goroutine exited.")
			wg.Done()
		}()
		for {
			for projectName := range data.Config.P {
				mySigningKey := []byte(data.Config.Internal.TG.BotToken)

				jwtExp := time.Duration(data.Config.Internal.JWTExpHours) * time.Hour
				expTime := time.Now().Add(jwtExp)
				claims := TokenClaims{
					projectName,
					jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(expTime),
						Issuer:    "cto",
					},
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				ss, _ := token.SignedString(mySigningKey)
				log.Println("Rotated auth token:", ss)

				data.SetObj(fmt.Sprintf("authToken-%s", projectName), ss, jwtExp)
			}
			<-tokenRotationTicker.C
		}
	}()

	go func() {
		defer data.CSend("All SLA checks exited.")
		wgSLA := new(sync.WaitGroup)
		wgSLA.Add(len(data.Config.P))

		for projectName := range data.Config.P {
			go func(pName string) {
				defer func() {
					data.CSend(fmt.Sprintf("SLA exited for project %s.", pName))
					wgSLA.Done()
				}()
				checking.SLAAggregator(data, pName)
			}(projectName)
		}

		wgSLA.Wait()
	}()

	go func() {
		handleSysSignals(data)
	}()

	data.CSend("I am up!")

	wg.Wait()
	data.CSend("All goroutines are done (or no one left alive). Main function will now exit.")
}

type TokenClaims struct {
	ProjectName string `json:"project_name"`
	jwt.RegisteredClaims
}
