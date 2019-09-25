package objects

import (
	"io"
	"log"
	"net/http"
	"strings"
)

const STORAGE_PATH = "./data"

func get(w http.ResponseWriter, r *http.Request) {

	obj:= strings.Split(r.URL.EscapedPath(), "/")[2]

	stream,e:=getStream(obj)

	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, e := io.Copy(w, stream); e != nil {
		log.Println(e)
	}
}
