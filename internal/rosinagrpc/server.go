package rosinagrpc

import (
	"context"
	"log"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func New() *RosinaGrpcServer {
	return &RosinaGrpcServer{jobs: make(chan Job)}
}

type RosinaGrpcServer struct {
	valvedprotos.UnimplementedRosinaSvcServer
	jobs chan Job
}

type Job struct {
	Duration time.Duration
}

// Open the valve
func (s *RosinaGrpcServer) Open(ctx context.Context, req *valvedprotos.OpenRequest) (*valvedprotos.OpenReply, error) {
	log.Printf("Adding job of duration")
	d := time.Second * time.Duration(req.Duration)
	s.jobs <- Job{Duration: d}

	return &valvedprotos.OpenReply{}, nil
}

func (s *RosinaGrpcServer) Jobs() <-chan Job {
	return s.jobs
}

func (s *RosinaGrpcServer) Close() {
	close(s.jobs)
}
