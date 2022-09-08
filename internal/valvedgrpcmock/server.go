package valvedgrpcmock

import (
	"context"
	"log"

	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func New() valvedprotos.ValvedSvcServer {
	return &valvedGrpcServer{status: false}
}

type valvedGrpcServer struct {
	valvedprotos.UnimplementedValvedSvcServer
	status bool
}

// Open the valve
func (s *valvedGrpcServer) Open(_ context.Context, _ *valvedprotos.OpenValveRequest) (*valvedprotos.OpenValveReply, error) {
	log.Printf("Valve open request received")
	s.status = true
	return &valvedprotos.OpenValveReply{Message: "Valve Opened"}, nil
}

// Close the valve
func (s *valvedGrpcServer) Close(_ context.Context, _ *valvedprotos.CloseValveRequest) (*valvedprotos.CloseValveReply, error) {
	log.Printf("Valve close request received")
	s.status = false
	return &valvedprotos.CloseValveReply{Message: "Valve Closed"}, nil
}

// Returns the status of the valve
func (s *valvedGrpcServer) Status(_ context.Context, _ *valvedprotos.StatusValveRequest) (*valvedprotos.StatusValveReply, error) {
	f := func() (valvedprotos.ValveStatus, error) {
		if s.status {
			return valvedprotos.ValveStatus_Open, nil
		}
		return valvedprotos.ValveStatus_Close, nil
	}

	st, err := f()
	return &valvedprotos.StatusValveReply{Status: st}, err
}
