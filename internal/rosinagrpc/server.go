package rosinagrpc

import (
	"context"
	"log"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func New() *RosinaGrpcServer {
	return &RosinaGrpcServer{
		aborts: make(chan struct{}),
		jobs:   make(chan Job),
	}
}

type RosinaGrpcServer struct {
	valvedprotos.UnimplementedRosinaSvcServer
	jobs   chan Job
	aborts chan struct{}
}

type Job struct {
	Duration time.Duration
}

// Open the valve
func (s *RosinaGrpcServer) OpenValve(ctx context.Context, req *valvedprotos.OpenRequest) (*valvedprotos.OpenReply, error) {
	log.Printf("Adding job of duration %d second", req.Duration)
	d := time.Second * time.Duration(req.Duration)
	s.jobs <- Job{Duration: d}

	return &valvedprotos.OpenReply{}, nil
}

// Close the valve
func (s *RosinaGrpcServer) CloseValve(ctx context.Context, req *valvedprotos.CloseRequest) (*valvedprotos.CloseReply, error) {
	s.aborts <- struct{}{}
	return &valvedprotos.CloseReply{}, nil
}

func (s *RosinaGrpcServer) Aborts() <-chan struct{} {
	return s.aborts
}

func (s *RosinaGrpcServer) Jobs() <-chan Job {
	return s.jobs
}

func (s *RosinaGrpcServer) Close() {
	close(s.jobs)
}
