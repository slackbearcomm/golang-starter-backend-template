package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Clients struct {
	PostgresConn *pgxpool.Pool
	AWSSession   *session.Session
	AWSRegion    string
	S3BucketName string
}
