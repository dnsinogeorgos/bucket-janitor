package verify

import (
	"fmt"
	"mime"
	"path/filepath"

	"github.com/dnsinogeorgos/bucket-janitor/internal/load"
	"github.com/vimeo/go-magic/magic"
)

// Object scans Key name and Bytes of object in order to verify type
func Object(object *load.Object) error {
	keyMime := mime.TypeByExtension(filepath.Ext(object.Key))
	dataMagic := magic.MimeFromBytes(object.Data)

	_, err := fmt.Printf("Scanning object %s\n", object.Key)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("MIME of file extension is: %s\n", keyMime)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("Magic type of bytes is: %s\n\n", dataMagic)
	if err != nil {
		return err
	}

	return nil
}
