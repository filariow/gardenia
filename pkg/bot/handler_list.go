package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	tele "gopkg.in/telebot.v3"
)

func (r *rosinaBot) listSkedule(c tele.Context) error {
	ss, err := r.SkedulerCli.ListSkedules(context.TODO(), &valvedprotos.ListSkedulesRequest{})
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
}
