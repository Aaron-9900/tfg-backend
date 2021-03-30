package aws

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Session struct {
	Session *s3.S3
}

func (s *S3Session) GetSignedUrl(filename string) (string, error) {
	req, _ := s.Session.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("tug-test"),
		Key:    aws.String(filename),
	})
	return req.Presign(5 * time.Minute)
}

func Init() *S3Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3")},
	)
	if err != nil {
		log.Println("Failed to sign request", err)
	}
	// Create S3 service client
	svc := s3.New(sess)
	return &S3Session{Session: svc}
}
