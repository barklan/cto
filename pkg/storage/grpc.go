package storage

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/barklan/cto/pkg/protos/porter"
	"google.golang.org/grpc"
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
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
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
	r, err := c.ProjectAlert(ctx, &pb.ProjectAlertRequest{
		Project: project,
		Message: message,
	})
	if err != nil {
		log.Printf("could not send grpc ProjectAlert: %v", err)
		return
	}
	log.Printf("grpc reply to ProjectAlert: %s", r.GetMessage())
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
	r, err := c.InternalAlert(ctx, &pb.Message{
		Message: message,
	})
	if err != nil {
		log.Printf("could not send grpc InternalAlert: %v", err)
		return
	}
	log.Printf("grpc reply to InternalAlert: %s", r.GetMessage())
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
	r, err := c.NewIssue(ctx, &pb.NewIssueRequest{
		Project:   projectID,
		Env:       env,
		Service:   service,
		Timestamp: timestamp,
		Key:       key,
		Flag:      flag,
	})
	if err != nil {
		log.Printf("could not send grpc NewIssue: %v", err)
		return
	}
	log.Printf("grpc reply to NewIssue: %s", r.GetMessage())
}
