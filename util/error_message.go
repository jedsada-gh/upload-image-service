package util

import (
	"fmt"
	"net/http"
)

// ErrorMessage is response message error request
func ErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprint(message)))
}
