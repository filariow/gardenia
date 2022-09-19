package rosina

import (
	"context"
	"log"
	"time"

	"github.com/filariow/gardenia/internal/rosinagrpc"
	"github.com/filariow/gardenia/pkg/valvedprotos"
)

func Skedule(ctx context.Context, cli valvedprotos.ValvedSvcClient, jobs <-chan rosinagrpc.Job) error {
	for job := range jobs {
		log.Println("Giving water to the garden")
		if _, err := cli.Open(ctx, &valvedprotos.OpenValveRequest{}); err != nil {
			log.Printf("Error giving water: %s", err)
			return err
		}

		log.Printf("Waiting for %d seconds: until %s UTC", job.Duration/time.Second, time.Now().UTC().Add(job.Duration))
		time.Sleep(job.Duration)

		log.Println("Stopping water to the garden")
		closeReq := valvedprotos.CloseValveRequest{}
		if _, err := cli.Close(context.TODO(), &closeReq); err != nil {
			log.Printf("Error giving water: %s", err)
			panic(err)
		}

		log.Println("Stopped water to the garden")
	}
	return nil
}
