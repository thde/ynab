package ynab_test

import (
	"context"
	"log"
	"time"

	"golang.org/x/time/rate"
	"thde.io/ynab"
)

func ExampleNew() {
	client, err := ynab.New(
		ynab.WithBearerToken("your-personal-access-token"),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.GetUserWithResponse(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_ = resp
}

func ExampleWithRateLimiter() {
	// Stay within the YNAB API limit of 200 requests per hour.
	limiter := rate.NewLimiter(rate.Every(time.Hour/200), 1)

	client, err := ynab.New(
		ynab.WithBearerToken("your-personal-access-token"),
		ynab.WithRateLimiter(limiter),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}
