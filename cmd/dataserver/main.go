package main

import (
	"github.com/huizluo/goss/app/dataserver/heartbeat"
	"github.com/huizluo/goss/app/dataserver/locate"
	"github.com/huizluo/goss/app/dataserver/objects"
	"github.com/huizluo/goss/app/dataserver/temp"
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
