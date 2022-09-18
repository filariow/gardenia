package rosinagrpc

import (
	"context"
	"log"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func New(c valvedprotos.ValvedSvcClient) valvedprotos.RosinaSvcServer {
	return &rosinaGrpcServer{cli: c}
}

type rosinaGrpcServer struct {
	valvedprotos.UnimplementedRosinaSvcServer
	cli valvedprotos.ValvedSvcClient
}

// Open the valve
func (s *rosinaGrpcServer) Open(ctx context.Context, req *valvedprotos.OpenRequest) (*valvedprotos.OpenReply, error) {
	log.Println("Giving water to the garden")
	if _, err := s.cli.Open(ctx, &valvedprotos.OpenValveRequest{}); err != nil {
		return nil, err
	}

	st := time.Duration(req.Duration) * time.Second
	log.Printf("Waiting for %d seconds: until %s UTC", req.Duration, time.Now().UTC().Add(st))
	time.Sleep(time.Second * time.Duration(req.Duration))

	log.Println("Stopping water to the garden")
	closeReq := valvedprotos.CloseValveRequest{}
	if _, err := s.cli.Close(ctx, &closeReq); err != nil {
		panic(err)
	}

	log.Println("Stopped water to the garden")

	return &valvedprotos.OpenReply{}, nil
}
