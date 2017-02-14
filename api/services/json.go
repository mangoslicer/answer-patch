package services

import (
	"encoding/json"
	"net/http"
)

//Marshals and writes JSON to the http.ResponseWriter
func PrintJSON(w http.ResponseWriter, content interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	postJSON, err := json.MarshalIndent(content, "", " ")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(postJSON)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

}
