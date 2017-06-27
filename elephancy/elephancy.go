package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"

	fe "github.com/poidl/elephancy/frontendserver"
	api "github.com/poidl/elephancy/restapi/go"
)

var frontendClientPath = "../frontendclient"
var pathFavicon = frontendClientPath + "/resources/favicon.ico"
var ftempl = frontendClientPath + "/resources/framenew.html"
var ftemplFingerpr = frontendClientPath + "/resources_fingerprinted/frame_fingerprinted.html"

var apiConf = api.Configuration{
	ApiAddress:     "127.0.0.1:8088",
	ApiDataFile:    "../restapi/data/pages.json",
	ApiResourceDir: "../restapi/resources"}

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

func staticcacheHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: catch direct requests to e.g. /staticcache/../
	// fmt.Println("Handler: static")
	// Aggressive caching. All files in /static must be fingerprinted for
	// cache-busting. Search for "google developers http caching"
	w.Header().Add("Cache-Control", "max-age=31536000")
	http.StripPrefix("/staticcache/", http.FileServer(http.Dir(frontendClientPath+"/resources_fingerprinted"))).ServeHTTP(w, r)
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

func makePagesHandler() func(w http.ResponseWriter, r *http.Request) {

	// check when template was last modified
	f, err := os.Open(ftemplFingerpr)
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	modtimeTemplate := fi.ModTime()

	return func(w http.ResponseWriter, r *http.Request) {

		page, err := fe.FindPageByKeyValue("prettyurl", r.URL.Path)
		// TODO: necessary?
		if page.Prettyurl != r.URL.Path {
			http.NotFound(w, r)
			return
		}

		pages, err := fe.ListPages()
		if err != nil {
			log.Fatal(err)
		}

		body, lastmodified, err := fe.GetPageContent(page.Id)
		if err != nil {
			log.Fatal(err)
		}

		// check if template is newer than content
		if modtimeTemplate.After(lastmodified) {
			lastmodified = modtimeTemplate
		}
		// send not modified if both template and content are older than
		// what's in the the client's cache
		if ifNotModifiedResponse(w, r, lastmodified) {
			return
		}

		templdat := make(map[string]interface{})
		templdat["Pages"] = pages
		templdat["Content"] = template.HTML(string(body))
		templdat["Metatitle"] = page.Metatitle

		// fill template
		templ, err := template.ParseFiles(ftemplFingerpr)

		templ.Execute(w, &templdat)
		w.Header().Add("Cache-Control", "no-cache")
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
}

// // makeFileServerHandler proxies fileserver on backend, provided that
// // the correct header is set (redirect if not)
// func makeFileServerHandler(rp *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ajax := r.Header.Get("myheader")
// 		if ajax == "XMLHttpRequest" {
// 			rp.ServeHTTP(w, r)
// 		} else {
// 			pages, err := fe.ListPages()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			page, err := pages.GetPageByLinksSelf(r.URL.Path)
// 			if err != nil {
// 				http.NotFound(w, r)
// 				return
// 			}
// 			http.Redirect(w, r, page.Prettyurl, 302)
// 		}
// 	}
// }

// // FileServerNew serves files from the frontend
// func FileServerNew(w http.ResponseWriter, r *http.Request) {
// 	http.FileServer(http.Dir("./")).ServeHTTP(w, r)
// }

func makeAPIHandler(rp *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rp.ServeHTTP(w, r)
	}
}

func backendCacheDefault(r *http.Response) error {
	r.Header.Add("Cache-Control", "no-cache")
	return nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go http.ListenAndServe(apiConf.ApiAddress, api.MyRouter(apiConf))
	time.Sleep(300 * time.Millisecond)

	backendURL, _ := url.Parse("http://" + apiConf.ApiAddress)
	frontendProxy := httputil.NewSingleHostReverseProxy(backendURL)
	frontendProxy.ModifyResponse = backendCacheDefault

	fe.SetupcacheNew()
	fe.GenerateFingerprintedTemplate(ftempl, ftemplFingerpr)

	http.HandleFunc("/", makePagesHandler())
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/staticcache/", staticcacheHandler)
	http.HandleFunc("/json/", jsonHandler)
	// http.HandleFunc("/fileserver/", makeFileServerHandler(frontendProxy))
	http.HandleFunc("/api/", makeAPIHandler(frontendProxy))
	http.ListenAndServe(":8080", nil)

}
