package aws

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageSession interface {
	GetPutSignedUrl(filename string) (string, error)
	GetGetSignedUrl(filename string) (string, error)
	GenerateFileName(filename string) (string, error)
}

type S3Session struct {
	Session *s3.S3
}

func (s *S3Session) GetPutSignedUrl(filename string) (string, error) {
	req, _ := s.Session.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("tug-test"),
		Key:    aws.String(filename),
	})
	return req.Presign(5 * time.Minute)
}
func (s *S3Session) GetGetSignedUrl(filename string) (string, error) {
	req, _ := s.Session.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("tug-test"),
		Key:    aws.String(filename),
	})
	return req.Presign(5 * time.Minute)
}
func (s *S3Session) GenerateFileName(filename string) (string, error) {
	extPosition := strings.LastIndex(filename, ".")
	fileExt := ""
	if extPosition > -1 {
		fileExt = filename[extPosition:]
	}
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?!-")
	b := make([]rune, 18)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	newName := filename + "_" + string(b) + fileExt
	return newName, nil
}

func Init() StorageSession {
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
