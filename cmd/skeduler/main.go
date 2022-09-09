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
	EnvAddress     = "ADDRESS"
	EnvApplication = "APPLICATION"
	DefaultAddress = "0.0.0.0:12000"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	return runServer()
}

func runServer() error {
	a := getAddress()
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
	app, err := getApplication()
	if err != nil {
		return nil, err
	}

	sk, err := skeduler.New(app)
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

func getApplication() (string, error) {
	if app := os.Getenv(EnvApplication); app != "" {
		return app, nil
	}

	return "", fmt.Errorf("Env var %s must be set", EnvApplication)
}
