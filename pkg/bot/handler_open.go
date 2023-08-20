package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	tele "gopkg.in/telebot.v3"
)

func (r *rosinaBot) openValve(c tele.Context) error {
	t, err := time.ParseDuration(c.Args()[0])
	if err != nil {
		return c.Send(fmt.Sprintf("Hey @filariow, sorry I did not understood how much time you want to give water (%s): %v", c.Args()[0], err))
	}

	d := uint64(t.Seconds())
	if _, err := r.RosinaCli.OpenValve(context.TODO(), &valvedprotos.OpenRequest{Duration: d}); err != nil {
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

			s, err := r.ValvedCli.Status(context.TODO(), &valvedprotos.StatusValveRequest{})
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

}
