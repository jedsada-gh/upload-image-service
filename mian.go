package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/upload-image-service/data"
	filetype "gopkg.in/h2non/filetype.v1"
)

var (
	httpPort      = os.Getenv("PORT")
	httpPortLocal = "8080"
	listenIP      = "localhost"
)

var (
	messageMethodNotAllowed  = "Method Not Allowed"
	messageFileNotSupported  = "File Not Support bacause type Image only"
	messageBucketNameInvalid = "Bucket Invalid"
	messageAPIKeyInvalid     = "API Key Invalid"
	messageRegionInvalid     = "Region Invalid"
	messageNoSuchFile        = "No such file"
)

const (
	methodPost = "POST"
	keyImage   = "file"
	keyBucket  = "bucket"
	keyRegion  = "region"
	keyAPIKey  = "api_key"
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
		validateValue(w, r)
	} else {
		errorMessage(w, http.StatusMethodNotAllowed, messageMethodNotAllowed)
	}
}

func validateValue(w http.ResponseWriter, r *http.Request) {
	bucket := r.FormValue(keyBucket)
	apiKey := r.FormValue(keyAPIKey)
	region := r.FormValue(keyRegion)
	if len(bucket) == 0 {
		errorMessage(w, http.StatusBadRequest, messageBucketNameInvalid)
	} else if len(apiKey) == 0 {
		errorMessage(w, http.StatusBadRequest, messageAPIKeyInvalid)
	} else if len(region) == 0 {
		errorMessage(w, http.StatusBadRequest, messageRegionInvalid)
	} else {
		fmt.Printf("%s\t%s\t%s\n", bucket, apiKey, region)
		file, _, err := r.FormFile(keyImage)
		if err != nil {
			errorMessage(w, http.StatusBadRequest, messageNoSuchFile)
		} else {
			defer file.Close()
			var model data.UploadImage
			model.APIKey = apiKey
			model.Bucket = bucket
			model.Region = region
			model.Image = file
			validateTypeFile(model, w, r)
		}
	}
}

func validateTypeFile(model data.UploadImage, w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(nil)
	_, error := io.Copy(buf, model.Image)
	if error != nil {
		errorMessage(w, http.StatusBadRequest, error.Error())
	} else if buf == nil {
		errorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	} else if filetype.IsImage(buf.Bytes()) {
		// TODO : upload file to aws s3
		w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", model.Image)))
	} else {
		errorMessage(w, http.StatusBadRequest, messageFileNotSupported)
	}
}
