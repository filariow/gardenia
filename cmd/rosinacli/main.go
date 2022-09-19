package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	EnvVarDurationInSec = "DURATION_IN_SEC"
	EnvVarRosinaAddress = "ROSINA_ADDRESS"
)

var (
	ErrMissingEnvVar = fmt.Errorf("missing env variable")
	ErrInvalidEnvVar = fmt.Errorf("invalid env variable content")
)

func main() {
	log.Println("Starting rosinacli")
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

	a, err := parseRosinaAddressFromEnv()
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

func parseRosinaAddressFromEnv() (*string, error) {
	a := os.Getenv(EnvVarRosinaAddress)
	if a == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingEnvVar, EnvVarRosinaAddress)
	}

	return &a, nil
}

func buildGrpcClient(address string) (valvedprotos.RosinaSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewRosinaSvcClient(conn)
	return cli, nil
}

func careGarden(ctx context.Context, cli valvedprotos.RosinaSvcClient, duration uint64) error {
	log.Println("Giving water to the garden")
	openReq := valvedprotos.OpenRequest{Duration: duration}
	if _, err := cli.Open(ctx, &openReq); err != nil {
		return err
	}
	return nil
}
