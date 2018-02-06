package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SuccessMessage is response message error request
func SuccessMessage(w http.ResponseWriter, obj interface{}) {
	model, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(model)
}
