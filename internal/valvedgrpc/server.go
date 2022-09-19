package valvedgrpc

import (
	"context"
	"fmt"
	"log"

	"github.com/filariow/gardenia/pkg/valve"
	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func New(d valve.Driver) *ValvedGrpcServer {
	return &ValvedGrpcServer{d: d}
}

type ValvedGrpcServer struct {
	valvedprotos.UnimplementedValvedSvcServer
	d valve.Driver

	openEvents  chan struct{}
	closeEvents chan struct{}
}

// Open the valve
func (s *ValvedGrpcServer) Open(_ context.Context, _ *valvedprotos.OpenValveRequest) (*valvedprotos.OpenValveReply, error) {
	log.Printf("Valve open request received")
	if err := s.d.SwitchOn(); err != nil {
		return nil, err
	}
	return &valvedprotos.OpenValveReply{Message: "Valve Opened"}, nil
}

// Close the valve
func (s *ValvedGrpcServer) Close(_ context.Context, _ *valvedprotos.CloseValveRequest) (*valvedprotos.CloseValveReply, error) {
	log.Printf("Valve open request received")
	if err := s.d.SwitchOff(); err != nil {
		return nil, err
	}
	return &valvedprotos.CloseValveReply{Message: "Valve Closed"}, nil
}

// Returns the status of the valve
func (s *ValvedGrpcServer) Status(_ context.Context, _ *valvedprotos.StatusValveRequest) (*valvedprotos.StatusValveReply, error) {
	ns := s.d.Status()
	f := func() (valvedprotos.ValveStatus, error) {
		switch ns {
		case valve.ValveOpen:
			return valvedprotos.ValveStatus_Open, nil
		case valve.ValveClose:
			return valvedprotos.ValveStatus_Close, nil
		}
		return valvedprotos.ValveStatus_Invalid, fmt.Errorf("Invalid valve status %v", ns)
	}

	st, err := f()
	return &valvedprotos.StatusValveReply{Status: st}, err
}

func (s *ValvedGrpcServer) OpenEvents() <-chan struct{} {
	return s.openEvents
}

func (s *ValvedGrpcServer) CloseEvents() <-chan struct{} {
	return s.closeEvents
}
