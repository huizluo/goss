package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"goss/app/dataserver/locate"
)

func put(w http.ResponseWriter, r *http.Request) {
	f, e := os.Create(os.Getenv("STORAGE_PATH") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
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

func del(w http.ResponseWriter, r *http.Request) {
	hash := strings.Split(r.URL.EscapedPath(), "/")[2]
	files, _ := filepath.Glob(os.Getenv("STORAGE_PATH") + "/objects/" + hash + ".*")
	if len(files) != 1 {
		return
	}
	locate.Del(hash)
	os.Rename(files[0], os.Getenv("STORAGE_PATH")+"/garbage/"+filepath.Base(files[0]))
}
