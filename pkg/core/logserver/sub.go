package logserver

import (
	"log"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/core/storage"
	"github.com/barklan/cto/pkg/loginput"
	"github.com/barklan/cto/pkg/rabbit"
)

func Subscriber(data *storage.Data, reqs chan<- loginput.LogRequest) {
	conn := rabbit.OpenMQ(data.Log)
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
			// if !data.VarExists(projectID, "") {
			// 	data.Log.Warn("rejecting log req", zap.String("project", projectID))
			// TODO add `continue` here after you made sure you have that flag
			// }
			reqs <- loginput.LogRequest{
				ProjectID: projectID,
				Body:      d.Body,
			}
			data.Log.Info("log req for added to local queue", zap.String("project", projectID))
		}
	}()

	data.Log.Info("sub is active")

	<-make(chan struct{})
}
