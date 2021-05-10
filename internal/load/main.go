package load

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

func Objects(AwsAccessKeyId, AwsSecretKey, S3Region string, S3Buckets []string) (map[string][]*Object, error) {
	// Create session
	session, err := NewSession(AwsAccessKeyId, AwsSecretKey, S3Region)
	if err != nil {
		return nil, err
	}

	// Create objects map and populate it
	s3ObjectMap := make(map[string][]*s3.Object)
	for _, bucket := range S3Buckets {
		objectCh := ListBucket(session, bucket)
		objectSlice := make([]*s3.Object, 0)

		_ = objectCh

		for object := range objectCh {
			objectSlice = append(objectSlice, object)
		}
		s3ObjectMap[bucket] = objectSlice
	}

	// Create downloader
	downloader := NewDownloader(session)

	// Retrieve object header data in map of values of type Object
	byteObjectSliceMap := make(map[string][]*Object)
	for bucket, s3Objects := range s3ObjectMap {
		byteObjectSlice := make([]*Object, 0)
		for _, s3Object := range s3Objects {
			byteObject, err := RetrieveObject(downloader, bucket, s3Object)
			if err != nil {
				return nil, err
			}

			byteObjectSlice = append(byteObjectSlice, byteObject)
		}
		byteObjectSliceMap[bucket] = byteObjectSlice
	}

	return byteObjectSliceMap, nil
}
