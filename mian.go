package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	filetype "gopkg.in/h2non/filetype.v1"
)

var (
	httpPort = os.Getenv("PORT")
	listenIP = "localhost"
	s3Region = ""
	s3Bucket = ""
)

var (
	messageMethodNotAllowed = "Method Not Allowed"
	messageFileNotSupported = "File Not Support bacause type Image only"
)

func errorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprint(message)))
}

func validateTypeFile(file multipart.File, w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(nil)
	_, error := io.Copy(buf, file)
	if error != nil {
		errorMessage(w, http.StatusBadRequest, error.Error())
	} else if buf == nil {
		errorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	}
	if filetype.IsImage(buf.Bytes()) {
		// TODO : upload file to aws s3
		fmt.Println("File is an image")
		w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", file)))
	} else {
		errorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	}
}

func handlerUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("file")
		if err != nil {
			errorMessage(w, http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()
		validateTypeFile(file, w, r)
	} else {
		errorMessage(w, http.StatusMethodNotAllowed, messageMethodNotAllowed)
	}
}

func main() {
	if httpPort == "" {
		httpPort = "8080"
	}
	http.HandleFunc("/upload", handlerUpload)
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
