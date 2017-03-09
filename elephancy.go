package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	fe "mystuff/elephancy/frontend"
	mj "mystuff/elephancy/json"
	sw "mystuff/elephancy/swagger"
	swBackend "mystuff/elephancy/swaggerbackend/go"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"
)

var pathFavicon = "staticcache/favicon.ico"
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

func pagesHandler(w http.ResponseWriter, r *http.Request) {
	templdat, modtime, err := mj.GetTemplateData(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	f, err := os.Open(ftemplFingerpr)
	defer f.Close()

	// check if  frame as been modified
	fi, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	modtimeTemplate := fi.ModTime()
	if modtimeTemplate.After(modtime) {
		modtime = modtimeTemplate
	}

	if ifNotModifiedResponse(w, r, modtime) {
		return
	}
	templ, err := template.ParseFiles(ftemplFingerpr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	templ.Execute(w, &templdat)
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

// func defaultHandler(w http.ResponseWriter, r *http.Request) {
// 	// if path is "/" or "/index.html"
// 	if rootPath.MatchString(r.URL.Path) {
// 		r.URL.Path = "/"
// 		pagesHandler(w, r)
// 		return
// 	}
// 	http.NotFound(w, r)
// 	return
// }

func makeHandleFunc(pages []mj.Page, page mj.Page) func(w http.ResponseWriter, r *http.Request) {

	// get the content
	contentaddr := backendIpPort + page.Links.Self
	resp, err := http.Get(contentaddr)
	if err != nil {
		log.Fatal(err)
	}

	// check when content was last modified
	lastmodified, _ := http.ParseTime(resp.Header.Get("Last-Modified"))

	// check when template was last modified
	f, err := os.Open(ftemplFingerpr)
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	modtimeTemplate := fi.ModTime()

	// check if template is newer than content
	if modtimeTemplate.After(lastmodified) {
		lastmodified = modtimeTemplate
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
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	templ.Execute(bw, &templdat)
	bw.Flush()
	return func(w http.ResponseWriter, r *http.Request) {
		if page.Prettyurl != r.URL.Path {
			http.NotFound(w, r)
			return
		}
		println("**************** contentaddr: " + contentaddr)
		// send not modified if both template and content are older
		if ifNotModifiedResponse(w, r, lastmodified) {
			return
		}
		w.Write(b.Bytes())
		w.Header().Add("Cache-Control", "no-cache")
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
}

func addHandleFuncs() {
	// Is it good to set up routing for Prettyurls? Perhaps better to
	// make a single handler and translate the Prettyurl to a page by
	// using the API function
	// page, err := sw.FindPageByPrettyURL(r.URL.Path)
	// Would be easier to handle created/deleted pages that way, since it's
	// necessary to (un-) register handler functions. But slower.
	pages, err := sw.ListPages()
	if err != nil {
		log.Fatal(err)
	}
	for _, page := range pages {
		println("ppppppppppppppppppppP: " + page.Prettyurl)
		http.HandleFunc(page.Prettyurl, makeHandleFunc(pages, page))
	}
}

func bla(r *http.Response) error {
	// Client must ask every time if there is a modified version.
	r.Header.Add("Cache-Control", "no-cache")
	return nil
}

func makeContentHandler(rp *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ajax := r.Header.Get("myheader")
		if ajax == "XMLHttpRequest" {
			rp.ServeHTTP(w, r)
		} else {
			// log.Fatal("This is broken")
			// fill in content
			println(r.URL.Path)
			page, err := sw.FindPageByPrettyURL(r.URL.Path)
			if err != nil {
				http.NotFound(w, r)
			}
			http.Redirect(w, r, page.Prettyurl, 302)
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	backendURL, _ := url.Parse(backendIpPort)
	go http.ListenAndServe(backendIpPort[len("http://"):], swBackend.MyRouter())
	time.Sleep(300 * time.Millisecond)

	// rpURL, err := url.Parse(backendServer.URL)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	frontendProxy := httputil.NewSingleHostReverseProxy(backendURL)
	// frontendProxy.ModifyResponse = bla
	fe.SetupcacheNew()
	fe.GenerateFingerprintedTemplate(ftempl, ftemplFingerpr)
	addHandleFuncs()
	// http.HandleFunc("/", pagesHandlerNew)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/frontend/staticcache/", staticcacheHandler)
	http.HandleFunc("/json/", jsonHandler)
	http.HandleFunc("/api/content/", makeContentHandler(frontendProxy))
	http.HandleFunc("/files/", filesHandler)
	http.ListenAndServe(":8080", nil)

	///////////////////////////////////////////
	// getCacheResources()
	// setupcacheNew()
	// // setupcache()
	// generateFingerprintedTemplate()
	// http.HandleFunc("/favicon.ico", faviconHandler)
	// http.HandleFunc("/frontend/staticcache/", staticcacheHandler)
	// http.HandleFunc("/", pagesHandler)
	// http.HandleFunc("/json/", jsonHandler)
	// http.HandleFunc("/content/", contentHandler)
	// http.HandleFunc("/files/", filesHandler)
	// http.ListenAndServe(":8080", nil)
}
