package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(model interface{}, w http.ResponseWriter, header int) {
	JSON, err := json.Marshal(model)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(header)
	w.Write(JSON)
}
