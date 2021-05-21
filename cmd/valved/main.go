package main

import (
	"log"
	"net"
	"os"

	"github.com/filariow/garden/internal/valvedgrpc"
	"github.com/filariow/garden/pkg/valve"
	"github.com/filariow/garden/pkg/valvedprotos"
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
	p1, p2 := os.Getenv("VPIN_1"), os.Getenv("VPIN_2")
	d := valve.New(p1, p2)

	s := grpc.NewServer()
	vs := valvedgrpc.New(d)
	valvedprotos.RegisterValvedSvcServer(s, vs)

	ls, err := net.Listen("tcp", ":12000")
	if err != nil {
		return err
	}
	return s.Serve(ls)
}
