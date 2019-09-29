package objects

import (
	"log"
	"net/http"
)

const (
	LISTEN_ADDRESS = "127.0.0.1:8060"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		log.Println("---put----")
		put(w, r)
		return
	}
	if m == http.MethodGet {
		log.Println("---get----")
		get(w, r)
		return
	}
	if m == http.MethodDelete {
		log.Println("---del----")
		del(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
