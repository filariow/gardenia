package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func parseArgs() (func(context.Context, *app) (*string, error), error) {
	if len(os.Args) == 1 {
		return nil, errors.New("no argument passed")
	}

	switch os.Args[1] {
	case "on":
		return SwitchOn, nil
	case "ON":
		return SwitchOn, nil
	case "1":
		return SwitchOn, nil
	case "off":
		return SwitchOff, nil
	case "OFF":
		return SwitchOff, nil
	case "0":
		return SwitchOff, nil
	}
	return nil, fmt.Errorf("Invalid argument: %s", os.Args[1])
}

func SwitchOn(ctx context.Context, a *app) (*string, error) {
	r, err := a.cli.Open(ctx, &valvedprotos.OpenValveRequest{})
	if err != nil {
		return nil, err
	}

	m := r.GetMessage()
	return &m, nil
}

func SwitchOff(ctx context.Context, a *app) (*string, error) {
	r, err := a.cli.Close(ctx, &valvedprotos.CloseValveRequest{})
	if err != nil {
		return nil, err
	}

	m := r.GetMessage()
	return &m, nil
}

func printErr(err error) {
	fmt.Println(err)
	fmt.Println(`

Usage: ./valve (on|ON|1|off|OFF|0)
	- on, ON, 1		open the valve
	- off, OFF, 0	close the valve`)
}

func run() error {
	f, err := parseArgs()
	if err != nil {
		printErr(err)
		return nil
	}

	address := os.Getenv("VALVED_ADDRESS")
	a := newApp(address)
	if err := a.setupClient(); err != nil {
		return err
	}
	defer a.Close()

	ctx := context.Background()
	s, err := f(ctx, a)
	if err != nil {
		return err
	}
	fmt.Printf(`Operation result: "%s"`, *s)
	return nil
}

type app struct {
	address string
	conn    *grpc.ClientConn
	cli     valvedprotos.ValvedSvcClient
}

func newApp(address string) *app {
	f := func() string {
		if address == "" {
			return ":12000"
		}
		return address
	}

	return &app{address: f()}
}

func (a *app) setupClient() error {
	conn, err := grpc.Dial(a.address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	a.conn = conn

	a.cli = valvedprotos.NewValvedSvcClient(conn)
	return nil
}

func (a *app) Close() {
	a.conn.Close()
}
