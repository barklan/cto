package main

import (
	"log"
	"time"

	"github.com/barklan/cto/pkg/loginput"
	postgres "github.com/barklan/cto/pkg/postgres"
	"github.com/jmoiron/sqlx"
)

func main() {
	log.Println("Starting...")

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
	defer rdb.Close()

	reqs := make(chan loginput.LogRequest, 5)

	go loginput.Serve(rdb, reqs)
	go loginput.Publisher(reqs)

	<-make(chan struct{})
}
