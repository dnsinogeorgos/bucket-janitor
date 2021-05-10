package load

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Object is a container for the header byteslice and the object key
type Object struct {
	Key  string
	Data []byte
}

// fetchHeaderBytes number of first bytes to retrieve for each object
const fetchHeaderBytes = 1024

// RetrieveObject returns object header bytes
func RetrieveObject(downloader *s3manager.Downloader, bucket string, s3Object *s3.Object) (*Object, error) {
	bytes := &aws.WriteAtBuffer{}
	_, err := downloader.Download(bytes, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(*s3Object.Key),
		Range:  aws.String("bytes=0-" + strconv.Itoa(fetchHeaderBytes-1)),
	})
	if err != nil {
		return &Object{}, err
	}

	byteObject := Object{
		Key:  *s3Object.Key,
		Data: bytes.Bytes(),
	}

	fmt.Printf("%s: %s\n", bucket, *s3Object.Key)

	return &byteObject, nil
}
