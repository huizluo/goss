package objects

import (
	"log"
	"net/http"
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
	if m == http.MethodPost {
		log.Println("-----post------")
		post(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
