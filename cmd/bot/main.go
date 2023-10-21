package main

import (
	"log"
	"os"
	"time"

	"github.com/filariow/gardenia/pkg/bot"
	tele "gopkg.in/telebot.v3"
)

func main() {
	log.Println("starting bot")
	pa, ok := os.LookupEnv("PROMETHEUS_ADDRESS")
	if !ok || pa == "" {
		log.Println("PROMETHEUS_ADDRESS env variable is not required but it is not defined or empty")
		os.Exit(1)
	}

	sa, ok := os.LookupEnv("SKEDULER_ADDRESS")
	if !ok || pa == "" {
		log.Println("SCHEDULER_ADDRESS env variable is not required but it is not defined or empty")
		os.Exit(1)
	}

	ra, ok := os.LookupEnv("ROSINA_ADDRESS")
	if !ok || ra == "" {
		log.Println("ROSINA_ADDRESS env variable is not required but it is not defined or empty")
		os.Exit(1)
	}

	va, ok := os.LookupEnv("VALVED_ADDRESS")
	if !ok || va == "" {
		log.Println("VALVED_ADDRESS env variable is not required but it is not defined or empty")
		os.Exit(1)
	}

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	cfg := bot.Config{
		BotSettings:       pref,
		AllowedIDs:        []int64{396136575},
		PrometheusAddress: pa,
		RosinaAddress:     ra,
		ValvedAddress:     va,
		SkedulerAddress:   sa,
	}

	b, err := bot.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("bot started")
	b.Start()
}
