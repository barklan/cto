package loginput

import (
	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/rabbit"
	"github.com/streadway/amqp"
)

type LogRequest struct {
	ProjectID string
	Body      []byte
}

func Publisher(reqs <-chan LogRequest) {
	defer log.Panicln("publisher exited")

	conn := rabbit.OpenMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	panicOnErr(err, "failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	panicOnErr(err, "failed to declare an exchange")

	for req := range reqs {
		err = ch.Publish(
			"logs", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				Headers: amqp.Table{
					"projectID": req.ProjectID,
				},
				ContentType: "text/plain",
				Body:        req.Body,
			})
		panicOnErr(err, "failed to publish a message")

		log.Printf("published loginginput for project %q", req.ProjectID)
	}
}

func panicOnErr(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
