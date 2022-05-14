package main

import (
	"log"
	"net"

	"github.com/filariow/gardenia/internal/schedulegrpc"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
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

	ls, err := net.Listen("tcp", "0.0.0.0:12001")
	if err != nil {
		return err
	}
	return s.Serve(ls)
}
