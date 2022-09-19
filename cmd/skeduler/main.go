package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/filariow/gardenia/internal/skedulegrpc"
	"github.com/filariow/gardenia/pkg/skeduler"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	EnvAddress       = "ADDRESS"
	EnvApplication   = "APPLICATION"
	EnvValvedImage   = "RUN_IMAGE"
	EnvRosinaAddress = "ROSINA_ADDRESS"

	DefaultAddress = "0.0.0.0:12000"
)

func main() {
	log.Println("Starting skeduler")
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	return runServer()
}

func runServer() error {
	a := getAddress()
	log.Printf("Starting server at %s", a)
	ls, err := net.Listen("tcp", a)
	if err != nil {
		return err
	}

	s, err := buildServer()
	if err != nil {
		return err
	}

	return s.Serve(ls)
}

func buildServer() (*grpc.Server, error) {
	app, err := getRequiredEnvVar(EnvApplication)
	if err != nil {
		return nil, err
	}

	vi, err := getRequiredEnvVar(EnvValvedImage)
	if err != nil {
		return nil, err
	}

	va := os.Getenv(EnvRosinaAddress)

	sk, err := skeduler.New(app, vi, &va)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	ss := skedulegrpc.New(sk)
	valvedprotos.RegisterSkeduleSvcServer(s, ss)

	return s, nil
}

func getAddress() string {
	if a := os.Getenv(EnvAddress); a != "" {
		return a
	}

	return DefaultAddress
}

func getRequiredEnvVar(env string) (string, error) {
	if app := os.Getenv(env); app != "" {
		return app, nil
	}

	return "", fmt.Errorf("Env var %s must be set", env)
}
