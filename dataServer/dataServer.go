package main

import (
	"go-storage/dataServer/heartbeat"
	"go-storage/dataServer/locate"
	"go-storage/dataServer/objects"
	"log"
	"net/http"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()

	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
