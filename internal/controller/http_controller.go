package controller

import (
	"fmt"
	"net/http"
)

var mux *http.ServeMux

func initMux() *http.ServeMux {

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
	return mux
}

func init() {
	mux = initMux()
}

func StartHTTPController() {

}

func getKeyValue(w http.ResponseWriter, req *http.Request) {
	var id string = req.PathValue("key")

	value, ok := memoryDB[id]
	if ok == false {
		fmt.Println("The value does not exist in the in memoryDB")
		return
	}
	fmt.Println(id + " : " + value)
}

func addKeyValue(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request accepted")
	var reqBody models.UserRequest
	var decoder = json.NewDecoder(req.Body)

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&reqBody)

	if err != nil {
		fmt.Printf("Error: Invalid JSON payload: %s", err)
	}

	_, ok := memoryDB[reqBody.Value]
	if ok == true {
		fmt.Println("Key already exists in the in memoryDB")
		return
	}

	memoryDB[reqBody.Key] = reqBody.Value
}

func putKeyValue(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("key")
	fmt.Println(id)
}

func deleteKeyValue(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("key")
	fmt.Println(id)
}
