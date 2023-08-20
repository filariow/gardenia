package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/filariow/gardenia/pkg/bot"
	tele "gopkg.in/telebot.v3"
)

func main() {
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
