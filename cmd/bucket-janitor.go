package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/service/s3"

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

	// Create session
	session, err := load.NewSession(c.AwsAccessKeyId, c.AwsSecretKey, c.S3Region)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	// Create objects map and populate it
	s3ObjectMap := make(map[string][]*s3.Object)
	for _, bucket := range c.S3Buckets {
		object, err := load.ListBucket(session, bucket)
		if err != nil {
			log.Printf("%s\n", err)
		}

		s3ObjectMap[bucket] = object
	}

	// Create downloader
	downloader := load.NewDownloader(session)

	// Retrieve object header data in map of values of type Object
	byteObjectSliceMap := make(map[string][]*load.Object)
	for bucket, s3Objects := range s3ObjectMap {
		byteObjectSlice := make([]*load.Object, 0)
		for _, s3Object := range s3Objects {
			byteObject, err := load.RetrieveObject(downloader, bucket, s3Object)
			if err != nil {
				log.Printf("%s\n", err)
			}

			byteObjectSlice = append(byteObjectSlice, byteObject)
		}
		byteObjectSliceMap[bucket] = byteObjectSlice
	}

	// Scan *Object with the magic library
	for _, byteObjects := range byteObjectSliceMap {
		for _, byteObject := range byteObjects {
			err = verify.ScanHeader(byteObject)
			if err != nil {
				log.Printf("%s\n", err)
			}
		}
	}
}
