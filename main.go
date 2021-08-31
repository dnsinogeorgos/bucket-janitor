package main

import (
	"flag"

	"github.com/dnsinogeorgos/bucket-janitor/internal/janitor"
)

func main() {
	configpath := flag.String("c", "bucket-janitor.json", "path to json config")
	flag.Parse()

	j := janitor.New(*configpath)
	j.Run()
}
