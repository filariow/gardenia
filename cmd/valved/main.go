package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/filariow/gardenia/internal/valvedgrpc"
	"github.com/filariow/gardenia/pkg/valve"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

const (
	DefaultUnixSocketAddr = "/tmp/valved.sock"
	EnvSocketAddr         = "VSOCKET_ADDR"
	EnvUnixSocketAddr     = "VSOCKET_ADDR_UNIX"
	EnvSwitchDuration     = "VSWITCH_DURATION"
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
	t, err := strconv.ParseUint(os.Getenv(EnvSwitchDuration), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing switch duration as uint64 (%s): %w", EnvSwitchDuration, err)
	}

	p1, p2 := os.Getenv("VPIN_1"), os.Getenv("VPIN_2")

	d := valve.New(p1, p2, t)

	s := grpc.NewServer()
	vs := valvedgrpc.New(d)
	valvedprotos.RegisterValvedSvcServer(s, vs)

	ls, err := listen()
	if err != nil {
		return err
	}

	return s.Serve(ls)
}

func listen() (net.Listener, error) {
	if a := getAddress(); a != "" {
		return net.Listen("tcp", a)
	}

	sa := getSocketAddr()
	if err := os.RemoveAll(sa); err != nil {
		log.Fatal(err)
	}

	return net.Listen("unix", sa)
}

func getAddress() string {
	return os.Getenv(EnvSocketAddr)
}

func getSocketAddr() string {
	a := os.Getenv(EnvUnixSocketAddr)
	if a == "" {
		log.Printf("using Default Socket Address (%s) because provided one (%s) is empty: '%s'", DefaultUnixSocketAddr, EnvUnixSocketAddr, a)
		return DefaultUnixSocketAddr
	}

	ss := strings.Split(a, string(os.PathSeparator))
	ss = ss[:len(ss)-1]
	d := strings.Join(ss, string(os.PathSeparator))
	if fi, err := os.Stat(d); os.IsNotExist(err) || !fi.IsDir() {
		log.Printf("using Default Socket Address (%s) because provided one (%s) is an invalid path: '%s'", DefaultUnixSocketAddr, EnvUnixSocketAddr, a)
		return DefaultUnixSocketAddr
	}

	return a
}
