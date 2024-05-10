package utils

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Body    map[string]any `json:"body"`
}

func SuccessfullyCreated(w http.ResponseWriter, code int, message string, body map[string]any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var response = &JsonResponse{Code: code, Message: message, Body: body}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SuccessfullyFetchedAll(w http.ResponseWriter, responseBody *JsonResponse) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func SuccessfullyFoundOne(w http.ResponseWriter, responseBody *JsonResponse) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusFound)

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
