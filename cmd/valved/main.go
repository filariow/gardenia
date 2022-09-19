package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/filariow/gardenia/internal/valvedgrpc"
	"github.com/filariow/gardenia/pkg/valve"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	addMetrics(vs.OpenEvents(), vs.CloseEvents())
	return s.Serve(ls)
}

func addMetrics(openEvents <-chan struct{}, closeEvents <-chan struct{}) {
	vs := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "valved_status",
		Help: "The status of the valved service",
	})

	go func() {
		for range openEvents {
			log.Printf("open events received")
			vs.Set(1)
		}
	}()

	go func() {
		for range closeEvents {
			log.Printf("close events received")
			vs.Set(0)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Printf("error starting metrics server")
		}
	}()
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
