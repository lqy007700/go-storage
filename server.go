package main

import (
	"go-storage/storage"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/objects/", storage.Handler)
	http.HandleFunc("/locate/", storage.Handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
