package main

import (
	"log"
	"net"

	"github.com/filariow/gardenia/internal/valvedgrpcmock"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	log.Println("Starting grpc server")
	s := grpc.NewServer()
	vs := valvedgrpcmock.New()
	valvedprotos.RegisterValvedSvcServer(s, vs)

	reflection.Register(s)

	ls, err := net.Listen("tcp", ":12000")
	if err != nil {
		return err
	}
	log.Println("Binded to port 12000")
	return s.Serve(ls)
}
