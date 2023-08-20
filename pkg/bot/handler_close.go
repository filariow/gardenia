package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	tele "gopkg.in/telebot.v3"
)

func (r *rosinaBot) closeValve(c tele.Context) error {
	if _, err := r.RosinaCli.CloseValve(context.TODO(), &valvedprotos.CloseRequest{}); err != nil {
		if serr := c.Send(fmt.Sprintf("I encountered an error closing the valve: %v", err)); serr != nil {
			log.Printf("error sending reply for valved error (%v) :%v", err, serr)
		}
		return err
	}

	return c.Send("I asked to close the valve, please check with /status if it has been closed correctly")
}
