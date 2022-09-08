package main

import (
	"log"
	"net"
	"os"

	"github.com/filariow/gardenia/internal/valvedgrpc"
	"github.com/filariow/gardenia/pkg/valve"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const SockAddr = "/tmp/valved.sock"

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

	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	ls, err := net.Listen("unix", SockAddr)
	if err != nil {
		return err
	}
	return s.Serve(ls)
}
