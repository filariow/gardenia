package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stianeikeland/go-rpio/v4"
)

const (
	EnvAddress = "ADDRESS"
	EnvPin     = "PIN"

	DefaultAddress = ":2113"
)

var (
	stats = atomic.Uint64{}

	flg = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "flow",
		Help: "Liters flowed",
	})

	flmg = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "flow_last_minute",
		Help: "Liters flowed in last minute",
	})
)

func addMetrics(address string) {
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				r := float64(stats.Swap(0)) / 30

				flg.Add(r)
				flmg.Set(r)

				log.Printf("%.3f Liters/min", r)
			}
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Printf("error starting metrics server")
		}
	}()
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	// read config
	p, err := readPinFromEnv()
	if err != nil {
		return err
	}

	a := readAddressFromEnvOrDefault()

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		return err
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	// create context
	ctx, cancel := context.WithCancel(context.Background())

	// configure input pin
	pin := rpio.Pin(p)
	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge) // enable falling edge event detection
	defer func() {
		pin.Detect(rpio.NoEdge) // disable edge event detection
		log.Printf("Pin Detect removed")
	}()

	// run metrics server
	addMetrics(a)

	// handle interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// wait for signal
		<-c

		// cancel context
		cancel()
	}()

	// run collector
	go runCollector(ctx, pin)

	// wait till context is valid
	<-ctx.Done()
	return nil
}

func runCollector(ctx context.Context, pin rpio.Pin) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if pin.EdgeDetected() {
				stats.Add(1)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Config

func readPinFromEnv() (int, error) {
	p := os.Getenv(EnvPin)
	if p == "" {
		return 0, fmt.Errorf(
			"Environment variable %s not set: it has to be an MCU (RPi BCM2835) Pinout", EnvPin)
	}

	pi, err := strconv.Atoi(p)
	if err != nil {
		return 0, fmt.Errorf(
			"Environment variable '%s' is invalid, it has to be an MCU (RPi BCM2835) Pinout, provided value is '%s': %w", EnvPin, p, err)
	}

	return pi, nil
}

func readAddressFromEnvOrDefault() string {
	a, ok := os.LookupEnv(EnvAddress)
	if !ok {
		return DefaultAddress
	}

	return a
}
