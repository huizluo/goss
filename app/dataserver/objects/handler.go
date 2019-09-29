package objects

import (
	"log"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(os.Getwd())
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
