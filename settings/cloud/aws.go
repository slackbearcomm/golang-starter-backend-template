package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewAWSSession(awsRegion, awsAccessKeyID, awsSecretAccessKey string) *session.Session {
	// return session.Must(session.NewSessionWithOptions(session.Options{}))

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})

	return sess
}
