package main

import (
	"go-storage/dataServer/heartbeat"
	"go-storage/dataServer/locate"
	"go-storage/dataServer/objects"
	"log"
	"net/http"
	"time"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()

	time.Sleep(2 * time.Second)

	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
