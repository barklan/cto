package logserver

import (
	"log"

	"github.com/barklan/cto/pkg/loginput"
	"github.com/barklan/cto/pkg/rabbit"
	"github.com/barklan/cto/pkg/storage"
)

func Subscriber(data *storage.Data, reqs chan<- loginput.LogRequest) {
	conn := rabbit.OpenMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicln("failed to open mq channel")
	}
	defer ch.Close()

	msgs := rabbit.OpenSubAutoAck(ch, "logs")

	go func() {
		// TODO should detect a situation where all cores reject request.
		// That means that projectID var not set in any core, but project exists in pg.
		for d := range msgs {
			projectID := d.Headers["projectID"].(string)
			if !data.VarExists(projectID, "") {
				log.Printf("rejecting log req for project %s\n", projectID)
				// TODO add `continue` here after you made sure you have that flag
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
