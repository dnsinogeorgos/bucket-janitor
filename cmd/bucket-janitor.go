package main

import (
	"context"
	"flag"
	"log"

	"github.com/dnsinogeorgos/bucket-janitor/internal/verify"

	ac "github.com/aws/aws-sdk-go-v2/config"

	"github.com/dnsinogeorgos/bucket-janitor/internal/config"
	"github.com/dnsinogeorgos/bucket-janitor/internal/load"
)

func main() {
	configpath := flag.String("c", "bucket-janitor.json", "path to json config")
	flag.Parse()

	// Load configuration
	c, err := config.NewConfig(*configpath)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	cfg, err := ac.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// var wg sync.WaitGroup

	headers := load.Headers(
		cfg,
		c.S3Buckets,
		c.Concurrency,
	)

	verify.Headers(headers, c.Concurrency)
}
