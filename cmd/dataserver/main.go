package main

import (
	"goss/app/dataserver/heartbeat"
	"goss/app/dataserver/locate"
	"goss/app/dataserver/objects"
	"goss/app/dataserver/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDR"), nil))
}
