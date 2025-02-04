package grpc

import (
	"flag"
	"fmt"
	"log"

	"github.com/ryvasa/go-super-farmer/pkg/env"
	pb "github.com/ryvasa/go-super-farmer/proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func InitGRPCClient(env *env.Env) (pb.ReportServiceClient, error) {
	reportServiceAddr := fmt.Sprintf("%s:%s", env.ReportService.Host, env.ReportService.Port)

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to report service: %v", err)
	}

	log.Printf("Connected to report service at %s", reportServiceAddr)
	return pb.NewReportServiceClient(conn), nil
}
