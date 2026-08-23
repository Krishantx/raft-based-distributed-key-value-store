package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	pb "raft-based-kv/proto"
)

func main() {
	fmt.Println("Hello World")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{key}", getKeyValue)
	mux.HandleFunc("POST /{key}", addKeyValue)
	mux.HandleFunc("PUT /{key}", putKeyValue)
	mux.HandleFunc("DELETE /{key}", deleteKeyValue)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Server Failed to Start: %s", err)
	}
	fmt.Println("Server Listning on Port: 8080")
}

func getKeyValue(w http.ResponseWriter, req *http.Request) {
	return
}

func addKeyValue(w http.ResponseWriter, req *http.Request) {
	return
}

func putKeyValue(w http.ResponseWriter, req *http.Request) {
	return
}

func deleteKeyValue(w http.ResponseWriter, req *http.Request) {
	return
}
