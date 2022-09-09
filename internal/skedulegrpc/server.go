package skedulegrpc

import (
	"context"

	"github.com/filariow/gardenia/pkg/skeduler"
	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func New(scheduler skeduler.Skeduler) valvedprotos.SkeduleSvcServer {
	return &skeduleGrpcServer{skeduler: scheduler}
}

type skeduleGrpcServer struct {
	valvedprotos.UnimplementedSkeduleSvcServer

	skeduler skeduler.Skeduler
}

func (s *skeduleGrpcServer) AddSkedule(ctx context.Context, req *valvedprotos.AddSkeduleRequest) (*valvedprotos.AddSkeduleReply, error) {
	cron := req.GetSkedule().GetCronSkedule()
	dur := req.GetSkedule().GetDurationSec()
	name, err := s.skeduler.AddJob(ctx, cron, uint64(dur))
	if err != nil {
		return nil, err
	}

	return &valvedprotos.AddSkeduleReply{JobName: name}, nil
}

func (s *skeduleGrpcServer) ListSkedules(ctx context.Context, req *valvedprotos.ListSkedulesRequest) (*valvedprotos.ListSkedulesReply, error) {
	jj, err := s.skeduler.ListJobs(ctx)
	if err != nil {
		return nil, err
	}

	ss := make([]*valvedprotos.Skedule, len(jj))
	for i, j := range jj {
		ss[i] = s.mapJobToSkedule(&j)
	}

	rep := valvedprotos.ListSkedulesReply{Skedules: ss}
	return &rep, nil
}

func (s *skeduleGrpcServer) GetSkedule(ctx context.Context, req *valvedprotos.GetSkeduleRequest) (*valvedprotos.GetSkeduleReply, error) {
	j, err := s.skeduler.GetJob(ctx, req.GetJobName())
	if err != nil {
		return nil, err
	}

	rep := valvedprotos.GetSkeduleReply{
		Skedule: s.mapJobToSkedule(j),
	}
	return &rep, nil
}

func (s *skeduleGrpcServer) DeleteSkedule(ctx context.Context, req *valvedprotos.DeleteSkeduleRequest) (*valvedprotos.DeleteSkeduleReply, error) {
	err := s.skeduler.RemoveJob(ctx, req.GetJobName())
	if err != nil {
		return nil, err
	}
	return &valvedprotos.DeleteSkeduleReply{}, nil
}

func (s *skeduleGrpcServer) mapJobToSkedule(j *skeduler.Job) *valvedprotos.Skedule {
	return &valvedprotos.Skedule{
		JobName:     j.JobName,
		CronSkedule: j.Schedule,
		DurationSec: int64(j.Duration),
	}
}
