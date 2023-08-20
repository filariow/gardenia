package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	tele "gopkg.in/telebot.v3"
)

func (r *rosinaBot) statusValve(c tele.Context) error {
	st, err := r.ValvedCli.Status(context.TODO(), &valvedprotos.StatusValveRequest{})
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
}
