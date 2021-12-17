package porter

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/barklan/cto/pkg/protos/porter"
	"google.golang.org/grpc"
)

var port = 50051

type server struct {
	pb.UnimplementedPorterServer
	data *Data
}

func (s *server) TelegramSend(ctx context.Context, in *pb.TelegramSendRequest) (*pb.TelegramSendReply, error) {
	projectID := in.GetProject()
	tgMessage := in.GetMessage()
	log.Printf("Received request to send tg message %q for project %q", tgMessage, projectID)
	go func() {
		s.data.PSend(projectID, tgMessage)
	}()

	return &pb.TelegramSendReply{Message: "ok"}, nil
}

func Serve(data *Data) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPorterServer(s, &server{data: data})
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
