package main

import (
	"flag"
	"log"

	"github.com/dnsinogeorgos/bucket-janitor/internal/config"
	"github.com/dnsinogeorgos/bucket-janitor/internal/load"
	"github.com/dnsinogeorgos/bucket-janitor/internal/verify"
)

func main() {
	configpath := flag.String("c", "bucket-janitor.json", "path to json config")
	flag.Parse()

	// Load configuration
	c, err := config.NewConfig(*configpath)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	byteObjectSliceMap, err := load.Objects(
		c.AwsAccessKeyId,
		c.AwsSecretKey,
		c.S3Region,
		c.S3Buckets,
	)

	// Scan *Object with the magic library
	for _, byteObjects := range byteObjectSliceMap {
		for _, byteObject := range byteObjects {
			err = verify.Object(byteObject)
			if err != nil {
				log.Printf("%s\n", err)
			}
		}
	}
}
