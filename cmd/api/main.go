package main

import (
	"context"
	"fmt"
	"smm-scrapper/util/scrapper"
)

func main() {
	go func() {

		defaultRepo := &scrapper.DefaultRepo{}

		telegramScrapper := scrapper.NewDefaultScrapper(defaultRepo)
		respch := make(chan struct{})
		ctx := context.Background()
		telegramScrapper.AuthorizeBot(ctx, respch)
		for {
			select {
			case <-respch:
				fmt.Println("Done authorizing")
			case <-ctx.Done():
				fmt.Println("Error occured")
			}
		}
		//go scrapper.RunAnalyticsTask(defaultRepo, &telegramScrapper)
	}()
	fmt.Println("Before infinite loop")
	for {

	}
	fmt.Println("Done running")
}
