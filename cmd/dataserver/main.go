package main

import (
	"goss/app/dataserver/heartbeat"
	"goss/app/dataserver/locate"
	"goss/app/dataserver/objects"
	"log"
	"net/http"
)

func main(){
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/",objects.Handler)
	log.Fatal(http.ListenAndServe(objects.LISTEN_ADDRESS,nil))
}
