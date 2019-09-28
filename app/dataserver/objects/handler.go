package objects

import "net/http"

const (
	STORAGE_PATH   = "../../data"
	LISTEN_ADDRESS = "127.0.0.1:8061"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
