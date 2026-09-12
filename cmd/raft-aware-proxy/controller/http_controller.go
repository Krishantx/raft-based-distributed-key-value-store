package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	service "raft-based-kv/cmd/raft-aware-proxy/service"
	"raft-based-kv/internal/models"
)

type HTTPController struct {
	proxyService service.ProxyService
	mux          *http.ServeMux
}

func NewHTTPController(proxyService service.ProxyService) *HTTPController {
	http_controller := HTTPController{
		proxyService: proxyService,
		mux:          http.NewServeMux(),
	}
	http_controller.mux.HandleFunc("GET /{key}", http_controller.getKeyValue)
	http_controller.mux.HandleFunc("POST /", http_controller.addKeyValue)
	http_controller.mux.HandleFunc("PUT /", http_controller.putKeyValue)
	http_controller.mux.HandleFunc("DELETE /{key}", http_controller.deleteKeyValue)

	return &http_controller
}

func (s *HTTPController) StartHTTPController() {
	err := http.ListenAndServe(":8080", s.mux)
	if err != nil {
		fmt.Printf("Server Failed to Start: %s", err)
	}
	fmt.Println("Server Listning on Port: 8080")
}

func (s *HTTPController) getKeyValue(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path[1:]
	keyValue := s.proxyService.GetKeyValue(key)
	json.NewEncoder(w).Encode(keyValue)
}

func (s *HTTPController) addKeyValue(w http.ResponseWriter, req *http.Request) {
	keyValue := s.proxyService.AddKeyValue()
	json.NewEncoder(w).Encode(keyValue)
}

func (s *HTTPController) putKeyValue(w http.ResponseWriter, req *http.Request) {
	var payload models.KeyValue

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err == nil {
		fmt.Printf("Error Decoding JSON: %s", err)
		http.Error(w, "Incorrect Payload", http.StatusBadRequest)
	}
	keyValue := s.proxyService.PutKeyValue()
	json.NewEncoder(w).Encode(keyValue)
}
func (s *HTTPController) deleteKeyValue(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path[1:]
	keyValue := s.proxyService.DeleteKeyValue(key)
	json.NewEncoder(w).Encode(keyValue)
}
