package bot

import (
	"log"
	"math/rand"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type RosinaBot interface {
	Start()
}

type rosinaBot struct {
	*tele.Bot

	config      Config
	RosinaCli   valvedprotos.RosinaSvcClient
	ValvedCli   valvedprotos.ValvedSvcClient
	SkedulerCli valvedprotos.SkeduleSvcClient
}

type Config struct {
	BotSettings       tele.Settings
	AllowedIDs        []int64
	PrometheusAddress string
	RosinaAddress     string
	ValvedAddress     string
	SkedulerAddress   string
}

func New(cfg Config) (RosinaBot, error) {
	rcli, err := buildRosinaGrpcClient(cfg.RosinaAddress)
	if err != nil {
		return nil, err
	}

	vcli, err := buildValvedGrpcClient(cfg.ValvedAddress)
	if err != nil {
		return nil, err
	}

	scli, err := buildSkedulerGrpcClient(cfg.SkedulerAddress)
	if err != nil {
		return nil, err
	}

	b, err := tele.NewBot(cfg.BotSettings)
	if err != nil {
		return nil, err
	}

	r := rosinaBot{
		Bot: b,

		config:      cfg,
		RosinaCli:   rcli,
		SkedulerCli: scli,
		ValvedCli:   vcli,
	}

	adminOnly := b.Group()
	fm := func(chats ...int64) tele.MiddlewareFunc {
		return func(next tele.HandlerFunc) tele.HandlerFunc {
			return func(c tele.Context) error {
				replies := []string{
					"a chi si figl?",
					"a chi appartien?",
					"chi si tu?",
					"a mammet!",
				}

				for _, chat := range chats {
					if chat != c.Sender().ID {
						i := rand.Intn(len(replies) - 1)
						err := c.Reply(replies[i])
						if err != nil {
							log.Printf("error replying: %s", err)
						}
						return err
					}
				}

				return nil
			}
		}
	}
	adminOnly.Use(fm(cfg.AllowedIDs...))
	adminOnly.Use(middleware.Whitelist(cfg.AllowedIDs...))

	adminOnly.Handle("/open", r.openValve)
	adminOnly.Handle("/close", r.closeValve)
	adminOnly.Handle("/status", r.statusValve)
	adminOnly.Handle("/plot", r.plot)
	adminOnly.Handle("/list", r.listSkedule)
	adminOnly.Handle("/add", r.addSkedule)
	adminOnly.Handle("/delete", r.deleteSkedule)

	return r, err
}
