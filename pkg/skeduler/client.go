package skeduler

import (
	"os"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	EnvSkedulerAddress           = "SKEDULER_SERVER_ADDRESS"
	DefaultSkedulerServerAddress = "localhost:12001"
)

func NewClientFromEnv() (valvedprotos.SkeduleSvcClient, error) {
	address := getSkedulerAddress()

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewSkeduleSvcClient(conn)
	return cli, nil
}

func getSkedulerAddress() string {
	if a := os.Getenv(EnvSkedulerAddress); a != "" {
		return a
	}

	return DefaultSkedulerServerAddress
}
