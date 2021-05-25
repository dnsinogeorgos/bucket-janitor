package verify

import (
	"fmt"
	"mime"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/dnsinogeorgos/bucket-janitor/internal/load"
	"github.com/vimeo/go-magic/magic"
)

type token struct{}

var types struct {
	sync.Mutex
	sync.Map
}

type typeSet struct {
	Mime  string
	Magic string
}

// Headers verifies all headers
func Headers(objectChannels map[string]chan *load.Header, limit int) {
	var wg sync.WaitGroup

	sem := make(chan token, limit)

	for bucket, headersCh := range objectChannels {
		wg.Add(1)
		go FanOutHeaders(bucket, headersCh, sem, &wg)
	}

	wg.Wait()

	fmt.Println()
	types.Range(PrintTypes)
}

func PrintTypes(key interface{}, value interface{}) bool {
	k := key.(typeSet)
	v := *value.(*uint64)

	fmt.Printf("%d\t%s\t\t%s\n", v, k.Mime, k.Magic)
	return true
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
	if header.Error != nil {
		fmt.Printf("\n%s %s: %s\n", bucket, header.Key, header.Error.Error())
		return
	}

	keyMime := mime.TypeByExtension(filepath.Ext(header.Key))
	if keyMime == "" {
		keyMime = "<no extension>"
	}

	dataMagic := magic.MimeFromBytes(header.Data)
	if dataMagic == "" {
		dataMagic = "<unknown mime>"
	}

	set := typeSet{
		keyMime,
		dataMagic,
	}
	types.Lock()
	if v, ok := types.Load(set); ok {
		atomic.AddUint64(v.(*uint64), 1)
	} else {
		v := new(uint64)
		*v++
		types.Store(set, v)
	}
	types.Unlock()

	_, _, _ = bucket, keyMime, dataMagic

	// fmt.Printf("%s: %s\n  * MIME:  %s\n  * Magic: %s\n", bucket, header.Key, keyMime, dataMagic)
	fmt.Printf(".")
}
