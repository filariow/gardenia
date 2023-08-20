package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	tele "gopkg.in/telebot.v3"
)

func (r *rosinaBot) addSkedule(c tele.Context) error {
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
	ns, err := r.SkedulerCli.AddSkedule(context.TODO(), &valvedprotos.AddSkeduleRequest{Skedule: &s})
	if err != nil {
		return c.Send(fmt.Sprintf("error adding job '%s - %ds': %v", ce, t, err))
	}
	return c.Send(fmt.Sprintf("job %s added", ns.GetJobName()))
}
