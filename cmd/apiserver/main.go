package main

import (
	"goss/app/apiserver/heartbeat"
	"goss/app/apiserver/locate"
	"goss/app/apiserver/objects"
	"goss/app/apiserver/temp"
	"goss/app/apiserver/versions"
	"log"
	"net/http"
)

func main() {
	go heartbeat.ListenHeartbeat()
	//对象操作
	http.HandleFunc("/objects/", objects.Handler)
	//
	http.HandleFunc("/temp/",temp.Handler)
	//对象定位
	http.HandleFunc("/locate/", locate.Handler)
	//对象版本
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(":8060", nil))
}
