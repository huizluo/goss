package main

import (
	"goss/app/apiserver/heartbeat"
	"goss/app/apiserver/locate"
	"goss/app/apiserver/objects"
	"log"
	"net/http"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/",locate.Handler)
	log.Fatal(http.ListenAndServe(":8060", nil))
}
