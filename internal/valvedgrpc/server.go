package valvedgrpc

import (
	"context"
	"fmt"

	"github.com/filariow/garden/pkg/valve"
	"github.com/filariow/garden/pkg/valvedprotos"
)

func New(d valve.Driver) valvedprotos.ValvedSvcServer {
	return &valvedGrpcServer{d: d}
}

type valvedGrpcServer struct {
	valvedprotos.UnimplementedValvedSvcServer
	d valve.Driver
}

// Open the valve
func (s *valvedGrpcServer) Open(_ context.Context, _ *valvedprotos.OpenValveRequest) (*valvedprotos.OpenValveReply, error) {
	if err := s.d.SwitchOn(); err != nil {
		return nil, err
	}
	return &valvedprotos.OpenValveReply{Message: "Valve Opened"}, nil
}

// Close the valve
func (s *valvedGrpcServer) Close(_ context.Context, _ *valvedprotos.CloseValveRequest) (*valvedprotos.CloseValveReply, error) {
	if err := s.d.SwitchOff(); err != nil {
		return nil, err
	}
	return &valvedprotos.CloseValveReply{Message: "Valve Closed"}, nil
}

// Returns the status of the valve
func (s *valvedGrpcServer) Status(_ context.Context, _ *valvedprotos.StatusValveRequest) (*valvedprotos.StatusValveReply, error) {
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
