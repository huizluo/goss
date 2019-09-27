package main

import (
	"go-oss/chapter6/apiServer/temp"
	"goss/app/dataserver/heartbeat"
	"goss/app/dataserver/locate"
	"goss/app/dataserver/objects"
	"log"
	"net/http"
)

func main(){
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/",objects.Handler)
	http.HandleFunc("/temp/",temp.Handler)
	log.Fatal(http.ListenAndServe(objects.LISTEN_ADDRESS,nil))
}
