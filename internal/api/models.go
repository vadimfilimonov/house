package api

import (
	"encoding/json"
	"net/http"
)

type BadRequestOuput struct {
	Message string `json:"message"`
}

func BadRequest(w *http.ResponseWriter, errorMessage string) {
	if w == nil {
		return
	}

	writer := (*w)

	response, err := json.Marshal(BadRequestOuput{Message: errorMessage})
	if err != nil {
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write(response)
}
