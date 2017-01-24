package swaggerbe

import (
	"net/http"
)

type Default struct {
}

func FindPageById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ListPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	http.ServeFile(w, r, "./json/pages.json")
}
