package bot

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promcommon "github.com/prometheus/common/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	tele "gopkg.in/telebot.v3"
)

const (
	dpi                  = 96
	plotDefaultStartTime = 30 * time.Minute
	plotDefaultEndTime   = 0 * time.Second
)

type plotInterval struct {
	start time.Duration
	end   time.Duration
}

func parsePlotArgs(args []string) (*plotInterval, error) {
	switch len(args) {
	case 0:
		return &plotInterval{
			start: plotDefaultStartTime,
			end:   plotDefaultEndTime,
		}, nil

	case 1:
		s, err := time.ParseDuration(args[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing start time: %w", err)
		}
		return &plotInterval{start: s, end: plotDefaultEndTime}, nil
	default:
		s, err := time.ParseDuration(args[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing start time: %v", err)
		}
		e, err := time.ParseDuration(args[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing end time: %v", err)
		}
		return &plotInterval{start: s, end: e}, nil
	}
}

func (r *rosinaBot) plot(c tele.Context) error {
	cfg := promapi.Config{Address: r.config.PrometheusAddress}
	pc, err := promapi.NewClient(cfg)
	if err != nil {
		return c.Send(fmt.Sprintf("Error creating the prometheus client: %v", err))
	}

  pi, err := parsePlotArgs(c.Args())
	if err != nil {
		return c.Send(err.Error())
	}

	a := prometheusv1.NewAPI(pc)
	st, et := time.Now().UTC().Add(-pi.start), time.Now().UTC().Add(-pi.end)
	scale := int(math.Ceil(float64(et.Second()-st.Second()) / 15))
	if scale == 0 {
		scale = 1
	}
	step := time.Duration(15*scale) * time.Second
	pr := prometheusv1.Range{
		Start: st,
		End:   et,
		Step:  step,
	}

	v, ww, err := a.QueryRange(context.TODO(), "flow_last_minute", pr)
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
