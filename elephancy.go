package main

import (
	"html/template"
	"io/ioutil"
	"log"
	api "mystuff/elephancy/api/go"
	fe "mystuff/elephancy/frontend"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"
)

var pathFavicon = "frontend/staticcache/resources/favicon.ico"
var backendIpPort = "http://127.0.0.1:8088"

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

var ftempl = "./frontend/templates/frame_new.html"
var ftemplFingerpr = "./frontend/templates/frame.html"

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

		// page, err := fe.FindPageByPrettyURL(r.URL.Path)
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

		// get the content
		contentaddr := backendIpPort + page.Links.Self
		resp, err := http.Get(contentaddr)
		if err != nil {
			log.Fatal(err)
		}

		// check when content was last modified
		lastmodified, _ := http.ParseTime(resp.Header.Get("Last-Modified"))

		// check if template is newer than content
		if modtimeTemplate.After(lastmodified) {
			lastmodified = modtimeTemplate
		}
		// send not modified if both template and content are older
		if ifNotModifiedResponse(w, r, lastmodified) {
			return
		}

		// read content
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
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

func makeContentHandler(rp *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ajax := r.Header.Get("myheader")
		if ajax == "XMLHttpRequest" {
			// why does adding the header here work? Could this be
			// overwritten by backend? If yes, use
			// frontendProxy.ModifyResponse = modrep
			// in the main function and
			// func modrep(r *http.Response) error {
			// 	r.Header.Add("Cache-Control", "no-cache")
			// 	return nil
			// }
			w.Header().Add("Cache-Control", "no-cache")
			rp.ServeHTTP(w, r)
		} else {
			page, err := fe.FindPageByKeyValue("linksself", r.URL.Path)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			http.Redirect(w, r, page.Prettyurl, 302)
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	backendURL, _ := url.Parse(backendIpPort)
	go http.ListenAndServe(backendIpPort[len("http://"):], api.MyRouter())
	time.Sleep(300 * time.Millisecond)

	frontendProxy := httputil.NewSingleHostReverseProxy(backendURL)

	fe.SetupcacheNew()
	fe.GenerateFingerprintedTemplate(ftempl, ftemplFingerpr)
	http.HandleFunc("/", makePagesHandler())
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/frontend/staticcache/", staticcacheHandler)
	http.HandleFunc("/json/", jsonHandler)
	http.HandleFunc("/api/content/", makeContentHandler(frontendProxy))
	http.HandleFunc("/files/", filesHandler)
	http.ListenAndServe(":8080", nil)

}
