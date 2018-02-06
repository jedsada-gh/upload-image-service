package utils

import (
	"fmt"
	"net/http"
)

// ErrorMessage exported function
func ErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprint(message)))
}
