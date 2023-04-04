package filestore

import (
	"bytes"
	"fmt"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FileStore struct {
	Session    *session.Session
	S3         *s3.S3
	BucketName string
	Region     string
}

func NewFilestore(sess *session.Session, region, bucketName string) *FileStore {
	return &FileStore{
		Session:    sess,
		S3:         s3.New(sess),
		Region:     region,
		BucketName: bucketName,
	}
}

func (fs *FileStore) ListBuckets() (*s3.ListBucketsOutput, *faulterr.FaultErr) {
	result, err := fs.S3.ListBuckets(nil)
	if err != nil {
		return nil, faulterr.NewInternalServerError(err.Error())
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	return result, nil
}

func (fs *FileStore) ListBucketItems() *faulterr.FaultErr {
	resp, err := fs.S3.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(fs.BucketName)})
	if err != nil {
		return faulterr.NewInternalServerError(err.Error())
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	return nil
}

func (fs *FileStore) UploadFile(filename string, fileObj []byte) (*dbmodels.File, *faulterr.FaultErr) {
	timestamp := time.Now().UTC().Unix()
	key := strconv.Itoa(int(timestamp)) + "_" + strings.ReplaceAll(filename, " ", "_")

	uploader := s3manager.NewUploader(fs.Session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &fs.BucketName,
		Key:    &key,
		Body:   bytes.NewReader(fileObj),
	})
	if err != nil {
		// return nil, faulterr.NewInternalServerError("file upload")
		return nil, faulterr.NewInternalServerError(err.Error())
	}

	obj := &dbmodels.File{
		Name: key,
		URL:  fs.generateURL(key),
	}

	return obj, nil
}

func (fs *FileStore) DownloadFile(key string) (*dbmodels.File, *faulterr.FaultErr) {
	downloader := s3manager.NewDownloader(fs.Session)

	file, err := os.Create(key)
	if err != nil {
		return nil, faulterr.NewInternalServerError(err.Error())
	}

	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("gogql-dev"),
			Key:    &key,
		})
	if err != nil {
		return nil, faulterr.NewInternalServerError(err.Error())
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	obj := &dbmodels.File{
		Name: key,
		URL:  fs.generateURL(key),
	}

	return obj, nil
}

func (fs *FileStore) generateURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", fs.BucketName, fs.Region, key)
}
