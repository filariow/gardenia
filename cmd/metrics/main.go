package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	EnvVarValvedAddress = "VALVED_ADDRESS"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	va := os.Getenv(EnvVarValvedAddress)
	if va == "" {
		return fmt.Errorf("valved address environment variable must be defined (%s)", EnvVarValvedAddress)
	}

	conn, err := grpc.Dial(va, grpc.WithInsecure())
	if err != nil {
		return err
	}
	cli := valvedprotos.NewValvedSvcClient(conn)

	vs := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "valved_status",
		Help: "The status of the valved service",
	})

	http.Handle("/metrics", func() http.Handler {
		rep, err := cli.Status(context.TODO(), &valvedprotos.StatusValveRequest{})
		if err != nil {
			log.Printf("error retrieving status from valve: %w", err)
			return promhttp.Handler()
		}

		if rep == nil {
			log.Printf("No reply received from valved")
			return promhttp.Handler()
		}

		s := func() float64 {
			if rep != nil {
				switch rep.GetStatus() {
				case valvedprotos.ValveStatus_Open:
					return 1
				case valvedprotos.ValveStatus_Close:
					return 0
				}
			}
			return -1
		}()

		vs.Set(s)
		log.Printf("Set status value to %f", s)

		return promhttp.Handler()
	}())

	return http.ListenAndServe(":2112", nil)
}
