package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	f, e := os.Create(STORAGE_PATH + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println("put obj")
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	io.Copy(f, r.Body)
}

func get(w http.ResponseWriter, r *http.Request) {

	f:=getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
 	if f==""{
 		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w,f)
}
