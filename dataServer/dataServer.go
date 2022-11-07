package main

import (
	"go-storage/dataServer/heartbeat"
	"go-storage/storage"
	"log"
	"net/http"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()

	http.HandleFunc("/objects/", storage.Handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
