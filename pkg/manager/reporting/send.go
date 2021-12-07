package reporting

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/barklan/cto/pkg/protos"
	"google.golang.org/grpc"
)

var (
	addr        = "ctopanel.com:50051"
	projectName = "nftg"
)

func GoReport(wg *sync.WaitGroup, message string) {
	defer wg.Done()
	if _, err := Report(message); err != nil {
		log.Printf("Failed to report back: %v", err)
	}
}

func Report(message string) (string, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Report(ctx, &pb.ReportRequest{
		Message:     message,
		ProjectName: projectName,
	})
	if err != nil {
		return "", err
	}
	reply := r.GetMessage()
	return reply, nil
}
