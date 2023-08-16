package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	ra, ok := os.LookupEnv("ROSINA_ADDRESS")
	if !ok || ra == "" {
		fmt.Printf("ROSINA_ADDRESS env variable is not required but it is not defined or empty")
		os.Exit(1)
	}

	va, ok := os.LookupEnv("VALVED_ADDRESS")
	if !ok || va == "" {
		fmt.Printf("VALVED_ADDRESS env variable is not required but it is not defined or empty")
		os.Exit(1)
	}

	rcli, err := buildRosinaGrpcClient(ra)
	if err != nil {
		fmt.Printf("error building the rosina grpc client: %s", err)
		os.Exit(1)
	}

	vcli, err := buildValvedGrpcClient(va)
	if err != nil {
		fmt.Printf("error building the valved grpc client: %s", err)
		os.Exit(1)
	}

	b.Handle("/help", func(c tele.Context) error {
		s := c.Sender()
		if s.ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		return c.Send(`The folowing commands are implemented:
/help
  prints this help message

/close
  closes the valve

/open DURATION_IN_SEC
  opens the valve for DURATION_IN_SEC seconds

/status
  prints the status of the valve`)
	})

	b.Handle("/open", func(c tele.Context) error {
		s := c.Sender()
		if s.ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		if len(c.Args()) == 0 {
			return c.Send("Hey @filariow, please specify an amount of seconds for the valve to stay open")
		}

		t, err := strconv.ParseUint(c.Args()[0], 10, 64)
		if err != nil {
			return c.Send(fmt.Printf("Hey @filariow, you provided an invalid amount of time: %v", err))
		}
		if t <= 0 {
			return c.Send(fmt.Printf("Hey @filariow, you provided an invalid amount of time, it needs to be greater than 0"))
		}

		if _, err := rcli.OpenValve(context.TODO(), &valvedprotos.OpenRequest{Duration: t}); err != nil {
			return c.Send(fmt.Sprintf("Something wrong happend while opening the valve: %s", err))
		}

		go func() {
			time.Sleep(time.Duration(t+1) * time.Second)

			isClosed := false
			for i := 0; i < 5; i++ {
				if isClosed {
					if err := c.Send("The valve is now closed"); err != nil {
						log.Printf("error sending reply: %v", err)
						continue
					} else {
						break
					}
				}

				s, err := vcli.Status(context.TODO(), &valvedprotos.StatusValveRequest{})
				if err != nil {
					if serr := c.Send(fmt.Sprintf("I encountered an error reading the valve status, can not tell if it is closed: %v", err)); serr != nil {
						log.Printf("error sending reply for valved error (%v) :%v", err, serr)
					}
				}
				isClosed = s.GetStatus() == valvedprotos.ValveStatus_Close
			}
			return
		}()

		for i := 0; i < 5; i++ {
			if err := c.Send(fmt.Sprintf("Hi @filariow, I opened the valve for you. It will close in %d seconds. I'll let you know", t)); err != nil {
				log.Printf("error sending reply: %v", err)
			} else {
				break
			}
		}

		return nil
	})

	b.Handle("/close", func(c tele.Context) error {
		s := c.Sender()
		if s.ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		if _, err := rcli.CloseValve(context.TODO(), &valvedprotos.CloseRequest{}); err != nil {
			if serr := c.Send(fmt.Sprintf("I encountered an error closing the valve: %v", err)); serr != nil {
				log.Printf("error sending reply for valved error (%v) :%v", err, serr)
			}
			return err
		}

		return c.Send("I asked to close the valve, please check with /status if it has been closed correctly")
	})

	b.Handle("/status", func(c tele.Context) error {
		s := c.Sender()
		if s.ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		st, err := vcli.Status(context.TODO(), &valvedprotos.StatusValveRequest{})
		if err != nil {
			if serr := c.Send(fmt.Sprintf("I encountered an error reading the valve status, can not tell if it is closed: %v", err)); serr != nil {
				log.Printf("error sending reply for valved error (%v) :%v", err, serr)
			}
			return err
		}
		stStr := func() string {
			switch st.GetStatus() {
			case valvedprotos.ValveStatus_Close:
				return "close"
			case valvedprotos.ValveStatus_Open:
				return "open"
			case valvedprotos.ValveStatus_Invalid:
				return "invalid"
			}
			return "unknown"
		}()

		return c.Send(fmt.Sprintf("The status of the valve is '%s'", stStr))
	})

	log.Printf("bot started")
	b.Start()
}

func buildValvedGrpcClient(address string) (valvedprotos.ValvedSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewValvedSvcClient(conn)
	return cli, nil
}

func buildRosinaGrpcClient(address string) (valvedprotos.RosinaSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewRosinaSvcClient(conn)
	return cli, nil
}
