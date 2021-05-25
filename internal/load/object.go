package load

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Header is a container for the header byteslice and the object key
type Header struct {
	Key   string
	Data  []byte
	Error error
}

// BundleHeaders retries the s3 object data headers asynchronously and immediately returns a map of *Header channels
func BundleHeaders(downloader *manager.Downloader, objectChannels map[string]<-chan types.Object, limit int) map[string]chan *Header {
	channels := make(map[string]chan *Header)
	sem := make(chan token, limit)

	for bucket, objectCh := range objectChannels {
		channels[bucket] = FunnelHeaders(downloader, bucket, objectCh, sem)
	}

	return channels
}

// FunnelHeaders retrieves headers from a channel and returns them in a slice
func FunnelHeaders(downloader *manager.Downloader, bucket string, objects <-chan types.Object, sem chan token) chan *Header {
	c := make(chan *Header)

	go func() {
		var wg sync.WaitGroup
		for object := range objects {
			wg.Add(1)
			sem <- token{}
			go func(object types.Object) {
				c <- RetrieveHeader(downloader, bucket, object)
				<-sem
				wg.Done()
			}(object)
		}

		wg.Wait()
		close(c)
	}()

	return c
}

// fetchHeaderBytes number of first bytes to retrieve for each object
const fetchHeaderBytes = 1024

// RetrieveHeader returns a byteslice of the object's data header
func RetrieveHeader(downloader *manager.Downloader, bucket string, object types.Object) *Header {
	b := make([]byte, fetchHeaderBytes)
	wb := manager.NewWriteAtBuffer(b)
	numBytes := fetchHeaderBytes

	switch {
	case numBytes > int(object.Size):
		numBytes = int(object.Size)
	case numBytes == 0:
		return &Header{
			Key:   *object.Key,
			Data:  make([]byte, 0),
			Error: fmt.Errorf("object has 0 bytes"),
		}
	}

	_, err := downloader.Download(context.Background(), wb, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(*object.Key),
		Range:  aws.String("bytes=0-" + strconv.Itoa(numBytes-1)),
	})
	if err != nil {
		return &Header{
			Key:   *object.Key,
			Data:  make([]byte, 0),
			Error: err,
		}
	}

	header := &Header{
		Key:   *object.Key,
		Data:  wb.Bytes(),
		Error: nil,
	}

	return header
}
