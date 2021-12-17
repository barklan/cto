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

func SendTgMessage(project, message string) {
	var addr string
	if v, ok := os.LookupEnv("CONFIG_ENV"); ok {
		if v == "dev" {
			addr = "localhost:" + port
		} else {
			addr = "cto_porter:" + port
		}
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewPorterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r, err := c.TelegramSend(ctx, &pb.TelegramSendRequest{
		Project: project,
		Message: message,
	})
	if err != nil {
		log.Printf("could not send grpc request to send tg message: %v", err)
		return
	}
	log.Printf("grpc reply to send tg message: %s", r.GetMessage())
}
