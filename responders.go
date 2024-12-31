package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		fmt.Println("Responding with 5xx error", msg)
	}

	respondWithJSON(w, code, NewApiFailureResponse(msg))
}

func respondWithJSON(w http.ResponseWriter, code int, payload ApiResponse) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal JSON response %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
