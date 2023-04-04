package config

import (
	"fmt"
	"gogql/settings/cloud"
	"gogql/settings/database/postgres"
	"log"
	"os"
	"strconv"

	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	defaultServerAddress string = ":8080"
)

// Config stores all configurations of the application
type Config struct {
	Server         *Server
	DBCreds        *DBCreds
	AWSCredentails *AWSCredentails
}

type Server struct {
	Address string
}

type DBCreds struct {
	Username string
	Password string
	Host     string
	Port     int64
	DBName   string
}

type AWSCredentails struct {
	Region          string
	AccessKeyID     string
	AccessKeySecret string
	S3BucketName    string
}

func LookupEnv(key string) (string, bool) {
	value, valid := os.LookupEnv(key)
	return strings.TrimSpace(value), valid
}

// This exists because os.GetEnv() would include bad characters like EOF in string values on some platforms (but not all)
func Getenv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	// Read env variables
	serverAddress := Getenv("SERVER_ADDRESS")
	awsRegion := Getenv("AWS_REGION")
	awsAccessKeyID := Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := Getenv("AWS_SECRET_ACCESS_KEY")
	awsS3BucketName := Getenv("AWS_S3_BUCKET_NAME")

	if serverAddress == "" {
		serverAddress = defaultServerAddress
		log.Println("WARNING: server address is missing, running on default port", defaultServerAddress)
	}

	if awsRegion == "" {
		log.Println("WARNING: aws default region is missing")
		//return nil, fmt.Errorf("aws default region is required")
	}
	if awsAccessKeyID == "" {
		log.Println("WARNING: aws access key id is missing")
		// return nil, fmt.Errorf("aws access key id is required")
	}
	if awsSecretAccessKey == "" {
		log.Println("WARNING: aws secret access key is missing")
		// return nil, fmt.Errorf("aws secret access key is required")
	}
	if awsS3BucketName == "" {
		log.Println("WARNING: aws s3 bucket name is missing")
		// return nil, fmt.Errorf("aws s3 bucket name is required")
	}

	awsCreds := &AWSCredentails{
		Region:          awsRegion,
		AccessKeyID:     awsAccessKeyID,
		AccessKeySecret: awsSecretAccessKey,
		S3BucketName:    awsS3BucketName,
	}

	dbCreds, err := fetchDbCreds(awsCreds)
	if err != nil {
		// No warning here, all configurations require dbCreds.
		log.Fatal(err)
	}
	server := &Server{
		Address: serverAddress,
	}
	config := &Config{
		DBCreds:        dbCreds,
		Server:         server,
		AWSCredentails: awsCreds,
	}

	return config, nil
}

func fetchDbCreds(aws_credentials *AWSCredentails) (*DBCreds, error) {
	AWS_DB_SECRET, AWS_DB_SECRET_EXISTS := LookupEnv("AWS_DB_SECRET")
	if AWS_DB_SECRET == "" {
		AWS_DB_SECRET_EXISTS = false
	}

	DB_USERNAME, DB_USERNAME_EXISTS := LookupEnv("DB_USERNAME")
	DB_PASSWORD, DB_PASSWORD_EXISTS := LookupEnv("DB_PASSWORD")
	DB_HOST, DB_HOST_EXISTS := LookupEnv("DB_HOST")
	DB_PORT_STR, DB_PORT_EXISTS := LookupEnv("DB_PORT")
	DB_NAME, DB_NAME_EXISTS := LookupEnv("DB_NAME")

	DB_PORT, err := strconv.ParseInt(DB_PORT_STR, 10, 64)
	if err != nil {
		DB_PORT_EXISTS = false
	}

	if DB_USERNAME_EXISTS && DB_PASSWORD_EXISTS && DB_HOST_EXISTS && DB_PORT_EXISTS && DB_NAME_EXISTS {
		log.Println("using database credentials defined in environment vars")

		return &DBCreds{
			Username: DB_USERNAME,
			Password: DB_PASSWORD,
			Host:     DB_HOST,
			Port:     DB_PORT,
			DBName:   DB_NAME,
		}, nil
	} else if AWS_DB_SECRET_EXISTS {
		log.Println("using database credentials defined in AWS secret")

		return GetAwsDBCreds(AWS_DB_SECRET, aws_credentials)
	} else {
		return nil, fmt.Errorf("no database credentials supplied")
	}
}

func GetAwsDBCreds(secretName string, awsCreds *AWSCredentails) (secret *DBCreds, err error) {
	// initiate aws session
	awsSession := cloud.NewAWSSession(awsCreds.Region, awsCreds.AccessKeyID, awsCreds.AccessKeySecret)

	svc := secretsmanager.New(awsSession,
		aws.NewConfig().WithRegion(awsCreds.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	// TODO: This is getting the following error
	// err.Error(): NoCredentialProviders: no valid providers in chain. Deprecated.

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				log.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return nil, err
	}

	var secretString, decodedBinarySecret string
	var dbSecret DBCreds

	if result.SecretString != nil {
		secretString = *result.SecretString
		err := json.Unmarshal([]byte(secretString), &dbSecret)
		if err != nil {
			log.Println("JSON Unmarshal Error:", err)
			return nil, err
		}
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Println("Base64 Decode Error:", err)
			return nil, err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])

		err1 := json.Unmarshal([]byte(decodedBinarySecret), &dbSecret)
		if err1 != nil {
			log.Println("JSON Unmarshal Error:", err)
			return nil, err
		}
	}

	return &dbSecret, nil
}

// Setup clients
func SetupClients(conf Config) *Clients {
	// initiate postgres connection
	psqlConn, err := postgres.ConnectPostgres(MakeDBSource(*conf.DBCreds))
	if err != nil {
		log.Fatal(err)
	}
	if psqlConn == nil {
		log.Fatal("unable to connect with postgres db")
	}

	// initiate aws session
	awsSession := cloud.NewAWSSession(conf.AWSCredentails.Region, conf.AWSCredentails.AccessKeyID, conf.AWSCredentails.AccessKeySecret)

	return &Clients{
		PostgresConn: psqlConn,
		AWSSession:   awsSession,
		AWSRegion:    conf.AWSCredentails.Region,
		S3BucketName: conf.AWSCredentails.S3BucketName,
	}
}

func MakeDBSource(dbCreds DBCreds) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbCreds.Host,
		dbCreds.Port,
		dbCreds.Username,
		dbCreds.Password,
		dbCreds.DBName,
	)
}
