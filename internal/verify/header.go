package verify

import (
	"fmt"
	"mime"
	"path/filepath"

	"github.com/dnsinogeorgos/bucket-janitor/internal/load"
	"github.com/vimeo/go-magic/magic"
)

// ScanHeader scans Key name and Bytes of object in order to verify type
func ScanHeader(header *load.Object) error {
	err := magic.AddMagicDir(magic.GetDefaultDir())
	if err != nil {
		return err
	}

	keyMime := mime.TypeByExtension(filepath.Ext(header.Key))
	dataMagic := magic.MimeFromBytes(header.Data)

	_, err = fmt.Printf("Scanning object %s\n", header.Key)
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
