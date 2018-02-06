package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/upload-image-service/utils"
	filetype "gopkg.in/h2non/filetype.v1"
)

var (
	httpPort      = os.Getenv("PORT")
	httpPortLocal = "8080"
	listenIP      = "localhost"
	s3Region      = ""
	s3Bucket      = ""
)

var (
	messageMethodNotAllowed = "Method Not Allowed"
	messageFileNotSupported = "File Not Support bacause type Image only"
)

const (
	methodPost = "POST"
	keyImage   = "file"
)

const (
	pathUpload = "/upload"
)

func main() {
	if httpPort == "" {
		httpPort = httpPortLocal
	}
	http.HandleFunc(pathUpload, handlerUpload)
	http.ListenAndServe(":"+httpPort, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))

	// s, err := session.NewSession(&aws.Config{Region: aws.String(s3Region)})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = addFileToS3(s, "result.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func handlerUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		file, _, err := r.FormFile(keyImage)
		if err != nil {
			utils.ErrorMessage(w, http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()
		validateTypeFile(file, w, r)
	} else {
		utils.ErrorMessage(w, http.StatusMethodNotAllowed, messageMethodNotAllowed)
	}
}

func validateTypeFile(file multipart.File, w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(nil)
	_, error := io.Copy(buf, file)
	if error != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, error.Error())
	} else if buf == nil {
		utils.ErrorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	} else if filetype.IsImage(buf.Bytes()) {
		// TODO : upload file to aws s3
		w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", file)))
	} else {
		utils.ErrorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	}
}

// func addFileToS3(s *session.Session, fileDir string) error {
// 	file, err := os.Open(fileDir)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	fileInfo, _ := file.Stat()
// 	var size = fileInfo.Size()
// 	buffer := make([]byte, size)
// 	file.Read(buffer)

// 	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
// 		Bucket:               aws.String(s3Bucket),
// 		Key:                  aws.String(fileDir),
// 		ACL:                  aws.String("private"),
// 		Body:                 bytes.NewReader(buffer),
// 		ContentLength:        aws.Int64(size),
// 		ContentType:          aws.String(http.DetectContentType(buffer)),
// 		ContentDisposition:   aws.String("attachment"),
// 		ServerSideEncryption: aws.String("AES256"),
// 	})
// 	return err
// }
