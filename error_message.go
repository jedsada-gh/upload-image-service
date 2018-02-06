package main

import (
	"fmt"
	"net/http"
)

// errorMessage is response message error request
func errorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprint(message)))
}
