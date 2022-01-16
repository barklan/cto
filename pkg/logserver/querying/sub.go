package querying

import (
	"encoding/json"

	"github.com/barklan/cto/pkg/porter"
	"github.com/barklan/cto/pkg/rabbit"
	"github.com/barklan/cto/pkg/storage"
	log "github.com/sirupsen/logrus"
)

func Subscriber(data *storage.Data, reqs chan<- porter.QueryRequest) {
	conn := rabbit.OpenMQ(data.Log)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicln("failed to open mq channel")
	}
	defer ch.Close()

	msgs := rabbit.OpenSubAutoAck(ch, "queries")

	go func() {
		// TODO should detect a situation where all cores reject request.
		// That means that projectID var not set in any core, but project exists in pg.
		for d := range msgs {
			projectID := d.Headers["projectID"].(string)
			if !data.VarExists(projectID, "") {
				log.WithField("project", projectID).Warn("rejecting log req ")
				// TODO add `continue` here after you made sure you have that flag
			}

			qr := porter.QueryRequest{}
			if err := json.Unmarshal(d.Body, &qr); err != nil {
				log.Panicln("failed to unmarshal query request from mq")
			}

			SetMsgInCache(data, qr.RequestID, porter.QWorking, "Query taken from rabbit by core replica.")

			reqs <- qr
			log.Printf("log req for project %s added to local queue", projectID)
		}
	}()

	data.Log.Info("sub is active")
	<-make(chan struct{})
}
