package load

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession creates new session for connections to AWS
func NewSession(awsAccessKeyId string, awsSecretKey string, s3Region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyId,
			awsSecretKey,
			""),
	})
	if err != nil {
		return &session.Session{}, err
	}

	return sess, nil
}
