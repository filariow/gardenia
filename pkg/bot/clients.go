package bot

import (
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

func buildValvedGrpcClient(address string) (valvedprotos.ValvedSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewValvedSvcClient(conn)
	return cli, nil
}

func buildRosinaGrpcClient(address string) (valvedprotos.RosinaSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewRosinaSvcClient(conn)
	return cli, nil
}

func buildSkedulerGrpcClient(address string) (valvedprotos.SkeduleSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewSkeduleSvcClient(conn)
	return cli, nil
}
