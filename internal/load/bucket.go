package load

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// batchSize is the number of objects to list on each API call to S3
const batchSize int = 1000

// ListBucket returns a slice of s3 objects
func ListBucket(sess *session.Session, bucket string) ([]*s3.Object, error) {
	client := s3.New(sess)
	responseLength := batchSize
	afterKey := ""
	batchSizeInt64 := int64(batchSize)
	objects := make([]*s3.Object, 0)

	fmt.Printf("listing bucket: %s\n", bucket)

	for responseLength == batchSize {
		response, err := client.ListObjectsV2(
			&s3.ListObjectsV2Input{
				Bucket:     aws.String(bucket),
				StartAfter: &afterKey,
				MaxKeys:    &batchSizeInt64,
			})
		if err != nil {
			return make([]*s3.Object, 0), err
		}

		responseLength = len(response.Contents)
		if responseLength == batchSize {
			afterKey = *response.Contents[responseLength-1].Key
		} else {
			afterKey = ""
		}

		objects = append(objects, response.Contents...)
	}

	return objects, nil
}
