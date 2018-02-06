package manager

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/upload-image-service/data"
)

var (
	awsAccessKey = os.Getenv("S3_ACCESS_KEY_PRIVATE")
	token        = ""
	pathImage    = "https://s3.amazonaws.com/api-upload-image/"
)

// UploadImageToS3 is upload image file to AWS S3
func UploadImageToS3(model data.UploadImage) (error, string) {
	fmt.Println(awsAccessKey)
	creds := credentials.NewStaticCredentials(awsAccessKey, model.APIKey, token)
	_, err := creds.Get()
	if err != nil {
		return err, ""
	}
	cfg := aws.NewConfig().WithRegion(model.Region).WithCredentials(creds)
	s, err := session.NewSession(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = addFileToS3(s, model)
	if err != nil {
		return err, ""
	}
	return nil, (pathImage + model.ImageName)
}

func addFileToS3(s *session.Session, model data.UploadImage) error {
	file, err := os.Create(model.ImageName)
	if err != nil {
		return err
	}

	_, err = file.Write(model.ImageByte)

	defer file.Close()
	if err != nil {
		return err
	}

	fileInfo, _ := file.Stat()
	fileName := fileInfo.Name()
	fileBytes := bytes.NewReader(model.ImageByte)
	fileType := http.DetectContentType(model.ImageByte)

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(model.Bucket),
		Key:           aws.String(fileName),
		Body:          fileBytes,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(fileType),
	})
	if err != nil {
		return err
	}
	err = os.Remove("./" + fileName)
	if err != nil {
		return err
	}
	return err
}
