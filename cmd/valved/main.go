package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/filariow/gardenia/internal/valvedgrpc"
	"github.com/filariow/gardenia/pkg/valve"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	DefaultSockAddr = "/tmp/valved.sock"
	EnvSocketAddr   = "VSOCKET_ADDR"
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

	sa := getSocketAddr()
	if err := os.RemoveAll(sa); err != nil {
		log.Fatal(err)
	}

	ls, err := net.Listen("unix", sa)
	if err != nil {
		return err
	}
	return s.Serve(ls)
}

func getSocketAddr() string {
	a := os.Getenv(EnvSocketAddr)
	if a == "" {
		log.Printf("using Default Socket Address (%s) because provided one (%s) is empty: '%s'", DefaultSockAddr, EnvSocketAddr, a)
		return DefaultSockAddr
	}

	ss := strings.Split(a, string(os.PathSeparator))
	ss = ss[:len(ss)-1]
	d := strings.Join(ss, string(os.PathSeparator))
	if fi, err := os.Stat(d); os.IsNotExist(err) || !fi.IsDir() {
		log.Printf("using Default Socket Address (%s) because provided one (%s) is an invalid path: '%s'", DefaultSockAddr, EnvSocketAddr, a)
		return DefaultSockAddr
	}
	return a
}
