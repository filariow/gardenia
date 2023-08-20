package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	tele "gopkg.in/telebot.v3"
)

func (r *rosinaBot) deleteSkedule(c tele.Context) error {
	if len(c.Args()) > 0 {
		j := c.Args()[0]
		if _, err := r.SkedulerCli.DeleteSkedule(context.TODO(), &valvedprotos.DeleteSkeduleRequest{JobName: j}); err != nil {
			return c.Send(fmt.Sprintf("error deleting %s: %v", j, err))
		}
		return c.Send(fmt.Sprintf("job %s deleted", j))
	}

	ss, err := r.SkedulerCli.ListSkedules(context.TODO(), &valvedprotos.ListSkedulesRequest{})
	if err != nil {
		return c.Send(fmt.Sprintf("error retrieving the schedules: %v", err))
	}

	rm := &tele.ReplyMarkup{
		ReplyKeyboard:   [][]tele.ReplyButton{{}},
		OneTimeKeyboard: true,
		Placeholder:     "/delete ",
	}
	var sb strings.Builder
	for _, s := range ss.Skedules {
		sb.WriteString(fmt.Sprintf("%s: %s - %s\n", s.GetJobName(), s.CronSkedule, (time.Duration(s.DurationSec) * time.Second).String()))
		t := fmt.Sprintf("/delete %s", s.GetJobName())
		b := tele.ReplyButton{Text: t}
		rm.ReplyKeyboard = append(rm.ReplyKeyboard, []tele.ReplyButton{b})
	}

	return c.Send(fmt.Sprintf("choose the one to delete:\n%s", sb.String()), r)
}
