package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/upload-image-service/data"
)

// ErrorMessage is response message error request
func ErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	var errorModel data.Error
	var errorDetailModel data.ErrorDetail
	errorDetailModel.Message = message
	errorModel.ErrorDetail = errorDetailModel
	error, err := json.Marshal(errorModel)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(error)
}
