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
		}

		if rep != nil {
			switch rep.GetStatus() {
			case valvedprotos.ValveStatus_Open:
				vs.Set(1)
			case valvedprotos.ValveStatus_Close:
				vs.Set(0)
			default:
				vs.Set(-1)
			}
		}

		return promhttp.Handler()
	}())

	return http.ListenAndServe(":2112", nil)
}
