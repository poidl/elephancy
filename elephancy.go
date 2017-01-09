package main

import (
	"html/template"
	"net/http"
	"os"
	"regexp"
	"time"
)

var pathFavicon = "staticcache/favicon.ico"

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

var templatefile = "./templ/frame.html"
var templates = template.Must(template.ParseFiles(templatefile))

func pagesHandler(w http.ResponseWriter, r *http.Request) {
	// // TODO: do this at start up
	// cachedat, err := loadJson("./frontend/staticcache/cache.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bla := cachedat["Frame.html"].(map[string]interface{})
	// println(bla["mystyle.css"].(string))
	// // err = templates.ExecuteTemplate(w, "frame.html", &cachedat)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	templdat, modtime, err := getTemplateData(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	f, err := os.Open(templatefile)
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// If frame as been modified after the template, set modtime to the more recent time. Note that template modifications which happen *after* the server has been started are accounted for in the Last-Modified header, but NOT DISPLAYED, perhaps because of template caching? As a consequence, if the template is modified while the server is running, and a user requests the page (gets the header, but not the updated content),than that user will NOT see the updated content if the server is restarted, because the browser has the latest timestamp and retrieves the old content from cache.
	modtimeTemplate := fi.ModTime()
	if modtimeTemplate.After(modtime) {
		modtime = modtimeTemplate
	}
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
	ajax := r.Header.Get("myheader")
	if ajax == "XMLHttpRequest" {
		http.FileServer(http.Dir("./")).ServeHTTP(w, r)
	} else {
		// fill in content
		urlpath, err := contentURLToUrlpath(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
		}
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
	setupcache()

	// 	http.HandleFunc("/favicon.ico", faviconHandler)
	// 	http.HandleFunc("/frontend/staticcache/", staticcacheHandler)
	// 	http.HandleFunc("/", pagesHandler)
	// 	http.HandleFunc("/json/", jsonHandler)
	// 	http.HandleFunc("/content/", contentHandler)
	// 	http.HandleFunc("/files/", filesHandler)
	// 	http.ListenAndServe(":8080", nil)
}
