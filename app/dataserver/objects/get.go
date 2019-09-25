package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {

	f, e := os.Open(STORAGE_PATH + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()

	if _, e := io.Copy(w, f); e != nil {
		log.Println(e)
	}
}
