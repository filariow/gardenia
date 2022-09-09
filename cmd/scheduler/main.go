package main

import (
	"log"
	"net"
	"os"

	"github.com/filariow/gardenia/internal/schedulegrpc"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	EnvAddress     = "ADDRESS"
	DefaultAddress = "0.0.0.0:12001"
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
	s := grpc.NewServer()
	ss := schedulegrpc.New()
	valvedprotos.RegisterScheduleSvcServer(s, ss)

	a := getAddress()
	ls, err := net.Listen("tcp", a)
	if err != nil {
		return err
	}
	return s.Serve(ls)
}

func getAddress() string {
	if a := os.Getenv(EnvAddress); a != "" {
		return a
	}

	return DefaultAddress
}
