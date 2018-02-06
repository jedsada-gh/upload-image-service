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
}

func handlerUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		file, _, err := r.FormFile(keyImage)
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

func validateTypeFile(file multipart.File, w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(nil)
	_, error := io.Copy(buf, file)
	if error != nil {
		errorMessage(w, http.StatusBadRequest, error.Error())
	} else if buf == nil {
		errorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	} else if filetype.IsImage(buf.Bytes()) {
		// TODO : upload file to aws s3
		w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", file)))
	} else {
		errorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	}
}
