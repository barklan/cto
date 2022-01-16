package porter

import (
	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/rabbit"
	"github.com/streadway/amqp"
)

type QueryRequest struct {
	RequestID string `json:"request_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	QueryText string `json:"query_text,omitempty"`
	Fields    string `json:"fields,omitempty"`
	Regex     string `json:"regex,omitempty"`
}

type QueryRequestWrap struct {
	ProjectID string
	QID       string
	Json      []byte
}

func Publisher(base *Base, queries <-chan QueryRequestWrap) {
	defer log.Panicln("publisher exited")

	conn := rabbit.OpenMQ(base.Log)
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
				Body: req.Json,
			})
		panicOnErr(err, "failed to publish a message")

		SetQRespInCache(base, req.QID, QWorking, "Query published to rabbit.")
		log.WithField("pid", req.ProjectID).Info("published query request")
	}
}

func panicOnErr(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", msg, err)
	}
}
