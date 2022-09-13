package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	EnvVarDurationInSec = "DURATION_IN_SEC"
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
	s, err := parseDurationInSecFromEnv()
	if err != nil {
		return err
	}
	log.Printf("Duration in Seconds: %d", *s)

	a, err := parseValvedAddressFromEnv()
	if err != nil {
		return err
	}
	log.Printf("Valved Address: %s", *a)

	cli, err := buildGrpcClient(*a)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := careGarden(ctx, cli, *s); err != nil {
		return err
	}

	return nil
}

func parseDurationInSecFromEnv() (*uint64, error) {
	d := os.Getenv(EnvVarDurationInSec)
	if d == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingEnvVar, EnvVarDurationInSec)
	}

	s, err := strconv.ParseUint(d, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEnvVar, EnvVarDurationInSec)
	}

	return &s, nil
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

func careGarden(ctx context.Context, cli valvedprotos.ValvedSvcClient, duration uint64) error {
	openReq := valvedprotos.OpenValveRequest{}
	if _, err := cli.Open(ctx, &openReq); err != nil {
		return err
	}
	log.Println("Giving water to the garden")

	st := time.Duration(duration) * time.Second
	log.Printf("Waiting for %d seconds: until %s UTC", duration, time.Now().UTC().Add(st))
	time.Sleep(st)

	closeReq := valvedprotos.CloseValveRequest{}
	if _, err := cli.Close(ctx, &closeReq); err != nil {
		panic(err)
	}
	log.Println("Stopped water to the garden")

	return nil
}
