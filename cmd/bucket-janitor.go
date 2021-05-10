package main

import (
	"flag"
	"log"
	"runtime"
	"sync"

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

	var wg sync.WaitGroup

	byteObjectSliceMap, err := load.Objects(
		c.AwsAccessKeyId,
		c.AwsSecretKey,
		c.S3Region,
		c.S3Buckets,
	)

	// Scan *Object with the magic library
	objectCh := make(chan *load.Object)
	for i := 0; i < runtime.NumCPU(); i++ {
		go verify.Object(objectCh, &wg)
	}
	for _, byteObjects := range byteObjectSliceMap {
		wg.Add(len(byteObjects))
		for _, byteObject := range byteObjects {
			objectCh <- byteObject
		}
	}
	wg.Wait()
}
