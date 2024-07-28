package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/filariow/gardenia/internal/rosinagrpc"
	"github.com/filariow/gardenia/pkg/rosina"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	EnvVarAddress       = "ADDRESS"
	EnvVarValvedAddress = "VALVED_ADDRESS"
)

var (
	ErrMissingEnvVar = fmt.Errorf("missing env variable")
	ErrInvalidEnvVar = fmt.Errorf("invalid env variable content")
)

func main() {
	log.Println("Starting rosina")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	va, err := parseValvedAddressFromEnv()
	if err != nil {
		return err
	}
	log.Printf("Valved Address: %s", *va)

	cli, err := buildGrpcClient(*va)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	rs := rosinagrpc.New()
	defer rs.Close()
	valvedprotos.RegisterRosinaSvcServer(s, rs)

	a, err := parseAddressFromEnv()
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", *a)
	if err != nil {
		return err
	}

	go rosina.Skedule(context.Background(), cli, rs.Jobs(), rs.Aborts())

	return s.Serve(lis)
}

func parseAddressFromEnv() (*string, error) {
	a := os.Getenv(EnvVarAddress)
	if a == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingEnvVar, EnvVarAddress)
	}

	return &a, nil
}

func parseValvedAddressFromEnv() (*string, error) {
	a := os.Getenv(EnvVarValvedAddress)
	if a == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingEnvVar, EnvVarValvedAddress)
	}

	return &a, nil
}

func buildGrpcClient(address string) (valvedprotos.ValvedSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewValvedSvcClient(conn)
	return cli, nil
}
