package verify

import (
	"fmt"
	"mime"
	"path/filepath"
	"sync"

	"github.com/dnsinogeorgos/bucket-janitor/internal/load"
	"github.com/vimeo/go-magic/magic"
)

// Object scans Key name and Bytes of object in order to verify type
func Object(objectCh chan *load.Object, wg *sync.WaitGroup) {
	for object := range objectCh {
		keyMime := mime.TypeByExtension(filepath.Ext(object.Key))
		dataMagic := magic.MimeFromBytes(object.Data)

		fmt.Printf("Scanning object %s\n", object.Key)
		fmt.Printf("MIME of file extension is: %s\n", keyMime)
		fmt.Printf("Magic type of bytes is: %s\n\n", dataMagic)

		wg.Done()
	}
}
