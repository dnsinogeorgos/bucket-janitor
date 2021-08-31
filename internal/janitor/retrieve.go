package janitor

import (
	"context"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Header struct {
	Bucket string
	Key    string
	Data   []byte
}

func (j *Janitor) retrieveHeaders() {
	defer j.wg.Done()
	defer close(j.headerChan)

	var wg sync.WaitGroup

	for object := range j.objectChan {
		wg.Add(1)
		j.headerSem <- token{}

		go j.retrieveHeader(object, &wg)
	}

	wg.Wait()
}

const fetchHeaderBytes = 1024

func (j *Janitor) retrieveHeader(o Object, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-j.headerSem }()

	b := make([]byte, fetchHeaderBytes)
	wb := manager.NewWriteAtBuffer(b)
	numBytes := fetchHeaderBytes

	switch {
	case numBytes > int(o.Size):
		numBytes = int(o.Size)
	case numBytes == 0:
		j.l.Info("object has 0 bytes", zap.String("bucket", o.Bucket), zap.String("key", o.Key))
		return
	}

	_, err := j.downloader.Download(context.Background(), wb, &s3.GetObjectInput{
		Bucket: aws.String(o.Bucket),
		Key:    aws.String(o.Key),
		Range:  aws.String("bytes=0-" + strconv.Itoa(numBytes-1)),
	})
	if err != nil {
		j.l.Error("error downloading header", zap.String("bucket", o.Bucket), zap.String("key", o.Key), zap.Error(err))
		return
	}

	j.headerChan <- Header{
		Key:  o.Key,
		Data: wb.Bytes(),
	}
}
