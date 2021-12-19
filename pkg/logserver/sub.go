package logserver

import (
	"log"

	"github.com/barklan/cto/pkg/loginput"
	"github.com/barklan/cto/pkg/rabbit"
	"github.com/barklan/cto/pkg/storage"
)

func panicOnErr(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}

func Subscriber(data *storage.Data, reqs chan<- loginput.LogRequest) {
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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	panicOnErr(err, "failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	panicOnErr(err, "failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	panicOnErr(err, "failed to register a consumer")

	go func() {
		// TODO should detect a situation where all cores reject request.
		// That means that projectID var not set in any core, but project exists in pg.
		for d := range msgs {
			projectID := d.Headers["projectID"].(string)
			if !data.VarExists(projectID, "") {
				log.Printf("rejecting log req for project %s\n", projectID)
			}
			reqs <- loginput.LogRequest{
				ProjectID: projectID,
				Body:      d.Body,
			}
			log.Printf("log req for project %s added to local queue", projectID)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-make(chan struct{})
}
