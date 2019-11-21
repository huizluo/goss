package main

import (
	"github.com/huizluo/goss/app/apiserver/heartbeat"
	"github.com/huizluo/goss/app/apiserver/locate"
	"github.com/huizluo/goss/app/apiserver/objects"
	"github.com/huizluo/goss/app/apiserver/temp"
	"github.com/huizluo/goss/app/apiserver/versions"
	"log"
	"net/http"
)

func main() {
	go heartbeat.ListenHeartbeat()
	//对象操作
	http.HandleFunc("/objects/", objects.Handler)
	//
	http.HandleFunc("/temp/", temp.Handler)
	//对象定位
	http.HandleFunc("/locate/", locate.Handler)
	//对象版本
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(":8060", nil))
}
