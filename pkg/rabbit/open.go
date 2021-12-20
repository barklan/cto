package rabbit

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/caarlos0/env"
	"github.com/streadway/amqp"
)

type MQConnectionData struct {
	Host     string `env:"RABBITMQ_HOST"`
	User     string `env:"RABBITMQ_DEFAULT_USER"`
	Password string `env:"RABBITMQ_DEFAULT_PASS"`
}

func OpenMQ() *amqp.Connection {
	cfg := MQConnectionData{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Panicln("failed to parse env for mq connection", err)
	}

	// TODO streadway/amqp is not maintained.
	// Migrate to official client eventually.
	// https://github.com/rabbitmq/amqp091-go
	var conn *amqp.Connection
	for i := 0; i < 30; i++ {
		conn, err = amqp.Dial(fmt.Sprintf(
			"amqp://%s:%s@%s:5672/",
			cfg.User,
			cfg.Password,
			cfg.Host,
		))
		if err != nil {
			if i == 29 {
				log.Panicln("failed to dial mq 30 times", err)
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	log.Println("connected to mq")
	return conn
}
