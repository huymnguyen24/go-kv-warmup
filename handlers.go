package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type API struct {
	store *Store
}

func (api *API) handleKey(writer http.ResponseWriter, request *http.Request) {
	key := strings.TrimPrefix(request.URL.Path, "/kv/")
	if key == "" {
		http.Error(writer, "Key is required", http.StatusBadRequest)
		return
	}
	switch request.Method {
	case http.MethodGet:
		value, ok := api.store.Get(key)
		if !ok {
			http.Error(writer, "Key not found", http.StatusNotFound)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(value))
	case http.MethodPut:
		request_body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		set_result := api.store.Set(key, string(request_body))
		if set_result {
			writer.WriteHeader(http.StatusCreated)
		} else {
			writer.WriteHeader(http.StatusOK)
		}
	case http.MethodDelete:
		delete_result := api.store.Delete(key)
		if !delete_result {
			http.Error(writer, "Key not found", http.StatusNotFound)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (api *API) handleList(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(api.store.Keys())
}
