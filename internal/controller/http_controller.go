package controller

//
import (
	"encoding/json"
	"fmt"
	"net/http"
	"raft-based-kv/internal/service"
)

var mux *http.ServeMux

type Controller struct {
	proxyService service.ProxyService
}

func NewHTTPController(proxyService service.ProxyService) *Controller {
	return &Controller{
		proxyService: proxyService,
	}
}

func (s *Controller) StartHTTPController(proxyService service.ProxyService) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{key}", s.getKeyValue)
	mux.HandleFunc("POST /{key}", s.addKeyValue)
	mux.HandleFunc("PUT /{key}", s.putKeyValue)
	mux.HandleFunc("DELETE /{key}", s.deleteKeyValue)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Printf("Server Failed to Start: %s", err)
	}
	fmt.Println("Server Listning on Port: 8080")
}

func (s *Controller) getKeyValue(w http.ResponseWriter, req *http.Request) {
	keyValue := s.proxyService.GetKeyValue()
	json.NewEncoder(w).Encode(keyValue)
}

func (s *Controller) addKeyValue(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("POST")
	keyValue := s.proxyService.AddKeyValue()
	json.NewEncoder(w).Encode(keyValue)
}

func (s *Controller) putKeyValue(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("PUT")
	keyValue := s.proxyService.PutKeyValue()
	json.NewEncoder(w).Encode(keyValue)
}
func (s *Controller) deleteKeyValue(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("DELETE")
	keyValue := s.proxyService.DeleteKeyValue()
	json.NewEncoder(w).Encode(keyValue)
}
