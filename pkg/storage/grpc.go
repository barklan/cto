package storage

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	pb "github.com/barklan/cto/pkg/protos/porter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port = "50051"

func (d *Data) dial() (*grpc.ClientConn, error) {
	var addr string
	if v, ok := os.LookupEnv("CONFIG_ENV"); ok {
		if v == "dev" {
			addr = "localhost:" + port
		} else {
			addr = "cto_porter:" + port
		}
	}
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	return conn, err
}

func (d *Data) ProjectAlert(project, message string) {
	conn, err := d.dial()
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewPorterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = c.ProjectAlert(ctx, &pb.ProjectAlertRequest{
		Project: project,
		Message: message,
	})
	if err != nil {
		log.WithError(err).Error("could not send grpc ProjectAlert")
		return
	}
}

func (d *Data) InternalAlert(message string) {
	conn, err := d.dial()
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewPorterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = c.InternalAlert(ctx, &pb.Message{
		Message: message,
	})
	if err != nil {
		log.WithError(err).Error("could not send grpc InternalAlert")
		return
	}
}

func (d *Data) NewIssue(projectID, env, service, timestamp, key, flag string) {
	conn, err := d.dial()
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewPorterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = c.NewIssue(ctx, &pb.NewIssueRequest{
		Project:   projectID,
		Env:       env,
		Service:   service,
		Timestamp: timestamp,
		Key:       key,
		Flag:      flag,
	})
	if err != nil {
		log.WithError(err).Error("could not send grpc NewIssue")
		return
	}
}
