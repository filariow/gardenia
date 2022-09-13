package main

import (
	"log"
	"net"
	"os"

	"github.com/filariow/gardenia/pkg/valvedproxy"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cs := os.Getenv("CONNECTION_STRING")
	p, err := valvedproxy.New(cs)
	if err != nil {
		return err
	}

	a := os.Getenv("ADDRESS")
	ls, err := net.Listen("tcp", a)
	if err != nil {
		return err
	}

	if err := p.Serve(ls); err != nil {
		return err
	}
	return nil
}
