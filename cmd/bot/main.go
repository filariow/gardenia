package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	promapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promcommon "github.com/prometheus/common/model"
	"google.golang.org/grpc"
	tele "gopkg.in/telebot.v3"
)

const dpi = 96

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

	sa, ok := os.LookupEnv("SKEDULER_ADDRESS")
	if !ok || pa == "" {
		fmt.Printf("SCHEDULER_ADDRESS env variable is not required but it is not defined or empty")
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

		t, err := time.ParseDuration(c.Args()[0])
		if err != nil {
			return c.Send(fmt.Sprintf("Hey @filariow, sorry I did not understood how much time you want to give water (%s): %v", c.Args()[0], err))
		}

		d := uint64(t.Seconds())
		if _, err := rcli.OpenValve(context.TODO(), &valvedprotos.OpenRequest{Duration: d}); err != nil {
			return c.Send(fmt.Sprintf("Something wrong happend while opening the valve: %s", err))
		}

		go func() {
			time.Sleep(t)

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
			if err := c.Send(fmt.Sprintf("Hi @filariow, I opened the valve for you. It will close in %d seconds. I'll let you know", d)); err != nil {
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
		bw := bufio.NewWriter(&b)
		d, err := plotData(m)
		if err != nil {
			return c.Send(fmt.Sprintf("error plotting the data: %v", err))
		}

		img := image.NewRGBA(image.Rect(0, 0, 10*dpi, 10*dpi))
		cv := vgimg.NewWith(vgimg.UseImage(img))
		d.Draw(draw.New(cv))
		if err := png.Encode(bw, cv.Image()); err != nil {
			return c.Send(fmt.Sprintf("error encoding the plot as png: %v", err))
		}

		if err := bw.Flush(); err != nil {
			return c.Send(fmt.Sprintf("error flushing the plot as png: %v", err))
		}

		br := bytes.NewReader(b.Bytes())
		p := &tele.Photo{File: tele.FromReader(br)}
		return c.Send(p)
	})

	cli, err := buildSkedulerGrpcClient(sa)
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/list", func(c tele.Context) error {
		if c.Sender().ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		ss, err := cli.ListSkedules(context.TODO(), &valvedprotos.ListSkedulesRequest{})
		if err != nil {
			return c.Send(fmt.Sprintf("error retrieving the schedules: %v", err))
		}

		var sb strings.Builder
		for _, s := range ss.Skedules {
			sb.WriteString(s.JobName)
			sb.WriteString(": ")
			sb.WriteString(s.CronSkedule)
			sb.WriteString(" - ")
			sb.WriteString((time.Duration(s.DurationSec) * time.Second).String())
			sb.WriteRune('\n')
		}
		return c.Send(sb.String())
	})

	b.Handle("/add", func(c tele.Context) error {
		if c.Sender().ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		if len(c.Args()) != 6 {
			return c.Send("six args expected")
		}

		aa := c.Args()
		ce := fmt.Sprintf("%s %s %s %s %s", aa[0], aa[1], aa[2], aa[3], aa[4])
		j := aa[5]

		t, err := time.ParseDuration(j)
		if err != nil || t < 0 {
			return c.Send(fmt.Sprintf("Invalid duration '%s': %v", j, err))
		}

		s := valvedprotos.Skedule{
			CronSkedule: ce,
			DurationSec: int64(t.Seconds()),
		}
		ns, err := cli.AddSkedule(context.TODO(), &valvedprotos.AddSkeduleRequest{Skedule: &s})
		if err != nil {
			return c.Send(fmt.Sprintf("error adding job '%s - %ds': %v", ce, t, err))
		}
		return c.Send(fmt.Sprintf("job %s added", ns.GetJobName()))
	})

	b.Handle("/delete", func(c tele.Context) error {
		if c.Sender().ID != 396136575 {
			return c.Send("Sorry, only @filariow is authorized")
		}

		if len(c.Args()) > 0 {
			j := c.Args()[0]
			if _, err := cli.DeleteSkedule(context.TODO(), &valvedprotos.DeleteSkeduleRequest{JobName: j}); err != nil {
				return c.Send(fmt.Sprintf("error deleting %s: %v", j, err))
			}
			return c.Send(fmt.Sprintf("job %s deleted", j))
		}

		ss, err := cli.ListSkedules(context.TODO(), &valvedprotos.ListSkedulesRequest{})
		if err != nil {
			return c.Send(fmt.Sprintf("error retrieving the schedules: %v", err))
		}

		r := &tele.ReplyMarkup{
			ReplyKeyboard:   [][]tele.ReplyButton{{}},
			OneTimeKeyboard: true,
			Placeholder:     "/delete ",
		}
		var sb strings.Builder
		for _, s := range ss.Skedules {
			sb.WriteString(fmt.Sprintf("%s: %s - %s\n", s.GetJobName(), s.CronSkedule, (time.Duration(s.DurationSec) * time.Second).String()))
			t := fmt.Sprintf("/delete %s", s.GetJobName())
			b := tele.ReplyButton{Text: t}
			r.ReplyKeyboard = append(r.ReplyKeyboard, []tele.ReplyButton{b})
		}

		return c.Send(fmt.Sprintf("choose the one to delete:\n%s", sb.String()), r)
	})

	log.Printf("bot started")
	b.Start()
}

func plotData(m promcommon.Matrix) (*plot.Plot, error) {
	p := plot.New()

	p.Title.Text = "Flow per minute"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Liters/min"

	xys := plotter.XYs{}
	if len(m) == 0 {
		return p, nil
	}

	v := m[0].Values[0]
	for _, a := range m[0].Values {
		xys = append(xys, plotter.XY{
			X: float64(a.Timestamp.Unix() - v.Timestamp.Unix()),
			Y: float64(a.Value),
		})
	}

	l, err := plotter.NewLine(xys)
	if err != nil {
		return nil, err
	}
	p.Add(l)

	return p, nil
}

func buildSkedulerGrpcClient(address string) (valvedprotos.SkeduleSvcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewSkeduleSvcClient(conn)
	return cli, nil
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
