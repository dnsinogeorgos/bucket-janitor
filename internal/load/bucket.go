package load

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// batchSize is the number of objects to list on each API call to S3
const batchSize int = 1000

// ListBucket returns a slice of s3 objects
func ListBucket(sess *session.Session, bucket string) chan *s3.Object {
	objectCh := make(chan *s3.Object)
	go func() {
		defer close(objectCh)

		client := s3.New(sess)
		responseLength := batchSize
		afterKey := ""
		batchSizeInt64 := int64(batchSize)

		fmt.Printf("start listing bucket: %s\n", bucket)

		for responseLength == batchSize {
			response, err := client.ListObjectsV2(
				&s3.ListObjectsV2Input{
					Bucket:     aws.String(bucket),
					StartAfter: &afterKey,
					MaxKeys:    &batchSizeInt64,
				})
			if err != nil {
				log.Printf("breaking loop listing bucket %s: %s\n", bucket, err)
				break
			}

			responseLength = len(response.Contents)
			if responseLength == batchSize {
				afterKey = *response.Contents[responseLength-1].Key
			} else {
				afterKey = ""
			}

			for _, object := range response.Contents {
				objectCh <- object
			}
		}

		fmt.Printf("end listing bucket: %s\n", bucket)
	}()
	return objectCh
}
