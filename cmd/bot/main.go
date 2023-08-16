package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	promapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promcommon "github.com/prometheus/common/model"
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

	pa, ok := os.LookupEnv("PROMETHEUS_ADDRESS")
	if !ok || pa == "" {
		fmt.Printf("PROMETHEUS_ADDRESS env variable is not required but it is not defined or empty")
		os.Exit(1)
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
		return c.Send(`The folowing commands are implemented:
/help
  prints this help message

/close
  closes the valve

/open DURATION_IN_SEC
  opens the valve for DURATION_IN_SEC seconds

/status
  prints the status of the valve

/plot
  plots the flow in a time frame`)
	})

	b.Handle("/open", func(c tele.Context) error {
		if c.Sender().ID != 396136575 {
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
		if c.Sender().ID != 396136575 {
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
		if c.Sender().ID != 396136575 {
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

	b.Handle("/plot", func(c tele.Context) error {
		if c.Sender().ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		cfg := promapi.Config{Address: pa}
		pc, err := promapi.NewClient(cfg)
		if err != nil {
			return c.Send(fmt.Sprintf("Error creating the prometheus client: %v", err))
		}

		s, err := time.ParseDuration(c.Args()[0])
		if err != nil {
			return c.Send(fmt.Sprintf("error parsing start time: %v", err))
		}
		e, err := time.ParseDuration(c.Args()[1])
		if err != nil {
			return c.Send(fmt.Sprintf("error parsing end time: %v", err))
		}

		a := prometheusv1.NewAPI(pc)
		r := prometheusv1.Range{
			Start: time.Now().UTC().Add(-s),
			End:   time.Now().UTC().Add(-e),
			Step:  15 * time.Second,
		}

		v, ww, err := a.QueryRange(context.TODO(), "flow_last_minute", r)
		if err != nil {
			return c.Send(fmt.Sprintf("Error retrieving flow_last_minute from prometheus: %v", err))
		}
		if len(ww) > 0 {
			log.Printf("warning querying range: %v", ww)
		}

		m := v.(promcommon.Matrix)
		var b bytes.Buffer
		foo := bufio.NewWriter(&b)
		d := plotData(m)
		cv := vgsvg.New(3*vg.Inch, 3*vg.Inch)
		d.Draw(draw.New(cv))
		if _, err := cv.WriteTo(foo); err != nil {
			return c.Send(fmt.Sprintf("error drawing the plot"))
		}
		foo.Flush()

		br := bytes.NewReader(b.Bytes())
		p := &tele.Photo{File: tele.FromReader(br)}
		return c.SendAlbum(tele.Album{p})
	})

	log.Printf("bot started")
	b.Start()
}

func plotData(m promcommon.Matrix) *plot.Plot {
	p := plot.New()

	p.Title.Text = "Flow per minute"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Liters/min"

	for _, e := range m {
		for _, a := range e.Values {
			plotutil.AddLinePoints(p, a.Timestamp.Time().Format(time.RFC1123), a.Value)
		}
	}

	return p
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
