package janitor

import (
	"mime"
	"path/filepath"
	"sync"

	"github.com/vimeo/go-magic/magic"
)

type TypeSet struct {
	Mime  string
	Magic string
}

func (j *Janitor) detectTypes() {
	defer j.wg.Done()
	defer close(j.typeSetChan)

	var wg sync.WaitGroup

	for h := range j.headerChan {
		wg.Add(1)
		go j.detectHeader(h, &wg)
	}

	wg.Wait()
}

func (j *Janitor) detectHeader(h Header, wg *sync.WaitGroup) {
	defer wg.Done()

	keyMime := mime.TypeByExtension(filepath.Ext(h.Key))
	if keyMime == "" {
		keyMime = "<no extension>"
	}

	dataMagic := magic.MimeFromBytes(h.Data)
	if dataMagic == "" {
		dataMagic = "<unknown mime>"
	}

	j.typeSetChan <- TypeSet{
		keyMime,
		dataMagic,
	}
}
