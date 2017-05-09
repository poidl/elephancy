/* OpenAPI spec version: 0.1.0
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * modified by S. Riha (2017)
 */

package api

import (
	"encoding/json"
	"log"
	mj "mystuff/elephancy/json"
	"net/http"
	"strconv"
)

var filename = "/home/stefan/programs/go/src/mystuff/elephancy/api/json/pages.json"

type Default struct {
}

// func FindPageById(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	// println("*************** Finding page by ID *****************")
// 	w.WriteHeader(http.StatusOK)
// }

func ListPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	http.ServeFile(w, r, filename)
}

// func FindPageByPrettyURL(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	filename := "/home/stefan/programs/go/src/mystuff/elephancy/json/pages.json"
// 	pages, err := mj.LoadPages(filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	v := r.URL.Query()
// 	// page, err := pages.PrettyURLToPage(v.Get("prettyurl"))
// 	page, err := pages.GetPageByPrettyURL(v.Get("prettyurl"))
// 	if err != nil {
// 		// println("notfound*****************************")
// 		http.NotFound(w, r)
// 	}
// 	json.NewEncoder(w).Encode(page)
// }

// func FindPageByLinksSelf(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	filename := "/home/stefan/programs/go/src/mystuff/elephancy/json/pages.json"
// 	pages, err := mj.LoadPages(filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	page, err := pages.GetPageByLinksSelf(r.URL.Path)
// 	println(page.Links.Self)
// 	if err != nil {
// 		// println("notfound*****************************")
// 		http.NotFound(w, r)
// 	}
// 	json.NewEncoder(w).Encode(page)
// }

// FindPageByKeyValue finds pages based on a single key-value pair of its properties
func FindPageByKeyValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pages, err := mj.LoadPages(filename)
	// println(pages[0].Prettyurl + "****************************")
	// println(pages[0].Links[0].Rel)
	if err != nil {
		log.Fatal(err)
	}
	vals := r.URL.Query()
	var page mj.Page
	switch vals["key"][0] {
	case "prettyurl":
		page, err = pages.GetPageByPrettyURL(vals["value"][0])

	default:
		http.NotFound(w, r)
	}
	// for k, v := range vals {
	// 	if k == "prettyurl" {
	// 		page, err = pages.GetPageByPrettyURL(v[0])
	// 	} else if k == "linksself" {
	// 		page, err = pages.GetPageByLinksSelf(v[0])
	// 	} else {
	// 		http.NotFound(w, r)
	// 	}
	// 	if err != nil {
	// 		http.NotFound(w, r)
	// 	}
	// }
	// if err != nil {
	// 	http.NotFound(w, r)
	// }
	// json.NewEncoder(os.Stdout).Encode(page)
	json.NewEncoder(w).Encode(page)
}

// Should be used by client-side only?
func GetPageContent(w http.ResponseWriter, r *http.Request) {
	ajax := r.Header.Get("myheader")
	if ajax == "XMLHttpRequest" {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")

		pages, err := mj.LoadPages(filename)
		if err != nil {
			log.Fatal(err)
		}
		ids := r.URL.Path[len("/api/content/"):]
		id, err := strconv.ParseInt(ids, 0, 64)
		if err != nil || (id == 0) {
			http.NotFound(w, r)
			return
		}
		page, err := pages.GetPageById(id)
		if err != nil {
			http.NotFound(w, r)
		}

		// get the content
		hrefSelf, err := page.GetLinkByRel("self")
		if err != nil {
			log.Fatal(err)
		}
		http.ServeFile(w, r, "."+hrefSelf)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Request is not an XMLHttpRequest."))
	}

}

// FileServer serves files WITHOUT caching policy
func FileServer(w http.ResponseWriter, r *http.Request) {
	// No caching policy here. Must be handled by frontend.
	http.StripPrefix("/fileserver/", http.FileServer(http.Dir("./fileserver"))).ServeHTTP(w, r)
}
