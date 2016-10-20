package main

import (
	"html/template"
	"net/http"
	"regexp"
	"time"
)

var pathFavicon string = "staticcache/favicon.ico"

// this is modified from package http
func ifNotModifiedResponse(w http.ResponseWriter, r *http.Request, modtime time.Time) bool {
	var unixEpochTime = time.Unix(0, 0)
	if modtime.IsZero() || modtime.Equal(unixEpochTime) {
		return false
	}
	if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		h := w.Header()
		delete(h, "Content-Type")
		delete(h, "Content-Length")
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	return false
}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")
var rootPath = regexp.MustCompile("^/$|^/index.html$")
var faviconPath = regexp.MustCompile("^/favicon.ico$")
var jsonPath = regexp.MustCompile("^/json/([a-zA-Z0-9]+).json$")
var contentPath = regexp.MustCompile("^/content/([a-zA-Z0-9]+).html$")

var templates = template.Must(template.ParseFiles("./templ/frame.html"))

func pagesHandler(w http.ResponseWriter, r *http.Request) {

	templdat, modtime, err := getTemplateData(r.URL.Path)
	// TODO: check if the frame as been modified and set modtime to the more recent time. Maybe no need to read frame.html since we use template caching. Is it possible to compare to the the time of compilation? Or maybe better time since server is started?
	if ifNotModifiedResponse(w, r, modtime) {
		return
	}
	err = templates.ExecuteTemplate(w, "frame.html", &templdat)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		return
	}
}

func staticcacheHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: catch direct requests to e.g. /staticcache/../
	// fmt.Println("Handler: static")
	// Aggressive caching. All files in /static must be fingerprinted for
	// cache-busting. Search for "google developers http caching"
	w.Header().Add("Cache-Control", "max-age=31536000")
	http.FileServer(http.Dir("./")).ServeHTTP(w, r)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	m := faviconPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, pathFavicon)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	m := jsonPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	http.FileServer(http.Dir("./")).ServeHTTP(w, r)
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./")).ServeHTTP(w, r)
}

func contentHandler(w http.ResponseWriter, r *http.Request) {
	m := contentPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	ajax := r.Header.Get("x-requested-with")
	if ajax == "XMLHttpRequest" {
		http.FileServer(http.Dir("./")).ServeHTTP(w, r)
	} else {
		// TODO: move this out of this file
		pcoll, err := loadPages()
		if err != nil {
			return
		}
		// fill in content
		contentURL := r.URL.Path
		pagename := pcoll.contentURLToPagename(contentURL)
		page := pcoll.getPage(pagename)
		urlpath := page["Urlpath"].(string)
		http.Redirect(w, r, urlpath, 302)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// if path is "/" or "/index.html"
	if rootPath.MatchString(r.URL.Path) {
		r.URL.Path = "/"
		pagesHandler(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func main() {
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/frontend/staticcache/", staticcacheHandler)
	http.HandleFunc("/", pagesHandler)
	http.HandleFunc("/json/", jsonHandler)
	http.HandleFunc("/content/", contentHandler)
	http.HandleFunc("/files/", filesHandler)
	http.ListenAndServe(":8080", nil)
}
