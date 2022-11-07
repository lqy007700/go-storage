package main

import (
	"go-storage/apiServer/heartbeat"
	"go-storage/apiServer/locate"
	"go-storage/apiServer/objects"
	"log"
	"net/http"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
