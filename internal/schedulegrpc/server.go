package schedulegrpc

import (
	"context"

	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func New() valvedprotos.ScheduleSvcServer {
	return &scheduleGrpcServer{}
}

type scheduleGrpcServer struct {
	valvedprotos.UnimplementedScheduleSvcServer
}

func (s *scheduleGrpcServer) AddSchedule(_ context.Context, _ *valvedprotos.AddScheduleRequest) (*valvedprotos.AddScheduleReply, error) {
	panic("not implemented") // TODO: Implement
}

func (s *scheduleGrpcServer) ListSchedules(_ context.Context, _ *valvedprotos.ListSchedulesRequest) (*valvedprotos.ListSchedulesReply, error) {
	panic("not implemented") // TODO: Implement
}

func (s *scheduleGrpcServer) GetSchedule(_ context.Context, _ *valvedprotos.GetScheduleRequest) (*valvedprotos.GetScheduleReply, error) {
	panic("not implemented") // TODO: Implement
}

func (s *scheduleGrpcServer) DeleteSchedule(_ context.Context, _ *valvedprotos.DeleteScheduleRequest) (*valvedprotos.DeleteScheduleReply, error) {
	panic("not implemented") // TODO: Implement
}
