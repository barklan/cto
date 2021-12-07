package grpcsrv

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/barklan/cto/pkg/protos"
	"github.com/barklan/cto/pkg/storage"
	"google.golang.org/grpc"
)

var port = 50051

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) Report(ctx context.Context, in *pb.ReportRequest) (*pb.ReportReply, error) {
	reportMessage := in.GetMessage()
	log.Printf("Received: %v", reportMessage)
	projectName := in.GetProjectName()

	storage.GData.PSend(projectName, reportMessage)
	return &pb.ReportReply{Message: "ok"}, nil
}

// TODO authentication!
func Serve(data *storage.Data) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		data.CSend(fmt.Sprintf("failed to listen: %v", err))
		log.Panicf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	data.CSend(fmt.Sprintf("grpc server listening at %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		data.CSend(fmt.Sprintf("failed to serve: %v", err))
		log.Panicf("failed to serve: %v", err)
	}
}
