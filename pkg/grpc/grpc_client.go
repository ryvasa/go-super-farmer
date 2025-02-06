package grpc

import (
	"flag"
	"fmt"
	"log"

	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	pb "github.com/ryvasa/go-super-farmer/proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGRPCClient(env *env.Env) (pb.ReportServiceClient, error) {
	var addr = flag.String("addr", env.ReportService.Host+":"+env.ReportService.Port, "the address to connect to")
	logrus.Log.Info(env.ReportService.Host + ":" + env.ReportService.Port)
	reportServiceAddr := fmt.Sprintf("%s:%s", env.ReportService.Host, env.ReportService.Port)

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to report service: %v", err)
	}

	log.Printf("Connected to report service at %s", reportServiceAddr)
	return pb.NewReportServiceClient(conn), nil
}
