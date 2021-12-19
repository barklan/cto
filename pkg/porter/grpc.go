package porter

import (
	"context"
	"fmt"
	"net"

	"github.com/barklan/cto/pkg/bot"
	pb "github.com/barklan/cto/pkg/protos/porter"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var port = 50051

type server struct {
	pb.UnimplementedPorterServer
	base  *Base
	sylon *bot.Sylon
}

func (s *server) ProjectAlert(ctx context.Context, in *pb.ProjectAlertRequest) (*pb.Message, error) {
	projectID := in.GetProject()
	tgMessage := in.GetMessage()
	log.Printf("Received request to send tg message %q for project %q", tgMessage, projectID)
	go func() {
		s.sylon.PSend(projectID, tgMessage)
	}()

	return &pb.Message{Message: "ok"}, nil
}

func (s *server) InternalAlert(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	message := in.GetMessage()
	s.sylon.CSend(message)
	return &pb.Message{Message: "ok"}, nil
}

func (s *server) NewIssue(ctx context.Context, in *pb.NewIssueRequest) (*pb.Message, error) {
	s.sylon.NotifyAboutError(
		in.GetProject(),
		in.GetEnv(),
		in.GetService(),
		in.GetTimestamp(),
		in.GetKey(),
		in.GetFlag(),
	)
	return &pb.Message{Message: "ok"}, nil
}

func (b *Base) ServeGRPC(sylon *bot.Sylon) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPorterServer(s, &server{
		base:  b,
		sylon: sylon,
	})

	log.Info("starting grpc server")
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
