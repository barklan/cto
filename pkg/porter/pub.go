package porter

import (
	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/rabbit"
	"github.com/streadway/amqp"
)

type QueryRequest struct {
	ProjectID string
	Text      []byte
}

func Publisher(queries <-chan QueryRequest) {
	defer log.Panicln("publisher exited")

	conn := rabbit.OpenMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	panicOnErr(err, "failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"queries", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	panicOnErr(err, "failed to declare an exchange")

	for req := range queries {
		err = ch.Publish(
			"queries", // exchange
			"",        // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				Headers: amqp.Table{
					"projectID": req.ProjectID,
				},
				ContentType: "text/plain",
				Body:        req.Text,
			})
		panicOnErr(err, "failed to publish a message")

		log.WithField("pid", req.ProjectID).Info("published query request")
	}
}

func panicOnErr(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
