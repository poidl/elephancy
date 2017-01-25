package swaggerbe

import (
	"encoding/json"
	"log"
	mj "mystuff/elephancy/json"
	"net/http"
)

type Default struct {
}

func FindPageById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	println("*************** Finding page by ID *****************")
	w.WriteHeader(http.StatusOK)
}

func ListPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	http.ServeFile(w, r, "./json/pages.json")
}

func FindPageByPrettyURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	filename := "/home/stefan/programs/go/src/mystuff/elephancy/json/pages.json"
	pcoll, err := mj.LoadJSONnew(filename)
	if err != nil {
		log.Fatal(err)
	}
	v := r.URL.Query()
	page, err := pcoll.PrettyURLToPage(v.Get("prettyurl"))
	if err != nil {
		http.NotFound(w, r)
	}
	json.NewEncoder(w).Encode(page)
}
