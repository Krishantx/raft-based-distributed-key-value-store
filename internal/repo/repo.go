package internal

import ()

var in_memory_db map[string]string

func init() {
	in_memory_db = make(map[string]string)
}

func getKeyValue(key string) string {
	return "Hello"
}

func addKeyValue(key string, value string) {

}

func deleteKeyValue(key string) {

}

func updateKeyValue(key string, value string) {

}
