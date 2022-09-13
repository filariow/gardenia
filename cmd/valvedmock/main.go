package main

import (
	"log"
	"net"
	"os"

	"github.com/filariow/gardenia/internal/valvedgrpcmock"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	log.Println("Starting grpc server")
	s := grpc.NewServer()
	vs := valvedgrpcmock.New()
	valvedprotos.RegisterValvedSvcServer(s, vs)

	reflection.Register(s)

	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	ls, err := net.Listen("unix", SockAddr)
	if err != nil {
		return err
	}
	log.Printf("Server started at %s", ls.Addr().String())
	return s.Serve(ls)
}
