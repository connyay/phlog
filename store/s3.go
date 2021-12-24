package store

import (
	"bytes"
	"io"
	"mime"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

var (
	_s3Bucket   = os.Getenv("AWS_S3_BUCKET")
	_s3Region   = os.Getenv("AWS_S3_REGION")
	_s3Key      = os.Getenv("AWS_S3_KEY")
	_s3Secret   = os.Getenv("AWS_S3_SECRET")
	_s3Endpoint = os.Getenv("AWS_S3_ENDPOINT")
)

type S3BlobStore struct {
	Store
}

func (S3BlobStore) AddBlob(data []byte, ext string) (string, error) {
	sess, err := connectAWS()
	if err != nil {
		return "", err
	}
	ref := uuid.NewString()
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(_s3Bucket),
		Key:         aws.String(ref),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(mime.TypeByExtension(ext)),
	})
	if err != nil {
		return "", err
	}
	return ref, nil
}

func (S3BlobStore) GetBlobByRef(ref string) (blob io.ReadCloser, ext string, err error) {
	sess, err := connectAWS()
	if err != nil {
		return nil, "", err
	}

	result, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(_s3Bucket),
		Key:    aws.String(ref),
	})
	if err != nil {
		return nil, "", err
	}
	if result.ContentType != nil {
		exts, _ := mime.ExtensionsByType(*result.ContentType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}

	return result.Body, ext, nil
}

func connectAWS() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(_s3Key, _s3Secret, ""),
		Endpoint:    aws.String(_s3Endpoint),
		Region:      aws.String(_s3Region),
	})
}
