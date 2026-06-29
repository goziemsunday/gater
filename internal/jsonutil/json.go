package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576 //1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		if err == io.EOF {
			return errors.New("request body is required")
		}
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteData(w http.ResponseWriter, status int, data any) {
	type envelope struct {
		Data any `json:"data"`
	}
	Write(w, status, &envelope{Data: data})
}

func WriteError(w http.ResponseWriter, status int, err any) {
	var errors []string

	switch e := err.(type) {
	case error:
		errors = append(errors, e.Error())
	case string:
		errors = append(errors, e)
	case []string:
		errors = e
	default:
		msg := fmt.Sprintf("%v", err)
		if msg == "<nil>" || msg == "" {
			msg = "an error occured"
		}
		errors = append(errors, msg)
	}

	type envelope struct {
		Error string `json:"error"`
	}

	Write(w, status, &envelope{Error: errors[0]})
}
