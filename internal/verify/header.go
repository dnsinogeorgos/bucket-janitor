package verify

import (
	"fmt"
	"mime"
	"path/filepath"
	"sync"

	"github.com/dnsinogeorgos/bucket-janitor/internal/load"
	"github.com/vimeo/go-magic/magic"
)

type token struct{}

// Headers verifies all headers
func Headers(objectChannels map[string]chan *load.Header, limit int) {
	var wg sync.WaitGroup

	sem := make(chan token, limit)

	for bucket, headersCh := range objectChannels {
		wg.Add(1)
		go FanOutHeaders(bucket, headersCh, sem, &wg)
	}

	wg.Wait()
}

// FanOutHeaders ships goroutines to verify headers
func FanOutHeaders(bucket string, headersCh chan *load.Header, sem chan token, wg *sync.WaitGroup) {
	for header := range headersCh {
		wg.Add(1)
		sem <- token{}
		go func(header *load.Header) {
			Header(bucket, header)
			<-sem
			wg.Done()
		}(header)
	}

	wg.Done()
}

// Header scans Key name and Bytes of object in order to verify type
func Header(bucket string, header *load.Header) {
	keyMime := mime.TypeByExtension(filepath.Ext(header.Key))
	dataMagic := magic.MimeFromBytes(header.Data)

	fmt.Printf("%s: %s\n  * MIME:  %s\n  * Magic: %s\n", bucket, header.Key, keyMime, dataMagic)
}
