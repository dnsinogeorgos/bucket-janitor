package load

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// NewDownloader creates a downloader object to use for downloading s3 objects
func NewDownloader(sess *session.Session) *s3manager.Downloader {
	downloader := s3manager.NewDownloader(sess)

	return downloader
}
