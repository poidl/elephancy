package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"
)

var pathFavicon string = "staticcache/favicon.ico"

func readContent(filename string) ([]byte, time.Time, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, time.Time{}, err //cannot use nil as type time.Time in return argument
	}
	// this is also done in ReadFile...inefficient
	// no need to check for error, since already done by ReadFile
	f, err := os.Open(filename)
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, time.Time{}, err
	}
	return content, fi.ModTime(), err
}

// loadPages opens a json file and returns the contents as a map[string]interface{}
// TODO: handle errors
func loadPages() (map[string]interface{}, error) {

	filename := "./json/pages.json"
	bytearr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var p interface{}
	err = json.Unmarshal(bytearr, &p)
	if err != nil {
		return nil, err
	}
	m := p.(map[string]interface{})
	return m, nil
}

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

// Urlpath2pagename returns the pagename of a requested url
func Urlpath2pagename(urlpath string, m map[string]interface{}) string {
	// TODO: return error in case it doesn't find anything
	for pagename, page := range m {
		pag := page.(map[string]interface{})
		if urlpath == pag["Urlpath"].(string) {
			return pagename
		}
	}
	return ""
}

//TODO: abstract this and Urlpath2pagename to one function

func Contenturl2pagename(contenturl string, m map[string]interface{}) string {
	// TODO: return error in case it doesn't find anything
	for pagename, page := range m {
		pag := page.(map[string]interface{})
		if contenturl == "/"+pag["ContentUrl"].(string) {
			return pagename
		}
	}
	return ""
}

// map2array orders page data into an array according linkweights. Itakes a map[string]interface{}, which maps pagenames to pages. A page is a json object containing page data as unordered list of key:value pairs. It returns an array of these objects, which is sorted corresponding to the "Linkweight" keys of the page objects.
func map2array(m map[string]interface{}) (arr []map[string]interface{}) {
	// we need an array for iteration order. see http://blog.golang.org/go-maps-in-action
	g := make(map[int]string)
	for pname, page := range m {
		pag := page.(map[string]interface{})
		g[int(pag["Linkweight"].(float64))] = pname
	}

	var keys []int
	for k := range g {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		arr = append(arr, m[g[k]].(map[string]interface{}))
	}
	return arr
}

func pagesHandler(w http.ResponseWriter, r *http.Request) {

	m, err := loadPages()
	if err != nil {
		return
	}

	// look up which page (pagename) is associated with the requested url
	pagename := Urlpath2pagename(r.URL.Path, m)
	if pagename == "" {
		http.NotFound(w, r)
		return
	}
	// load the page data into a map[string]interface{}
	page := m[pagename].(map[string]interface{})
	// read content from html file and get modification time
	content, modtime, err := readContent(page["ContentUrl"].(string))

	// TODO: check if the frame as been modified and set modtime to the more recent time. Maybe no need to read frame.html since we use template caching. Is it possible to compare to the the time of compilation? Or maybe better time since server is started?
	if ifNotModifiedResponse(w, r, modtime) {
		return
	}

	arr := map2array(m)
	blab := make(map[string]interface{})
	blab["Pages"] = arr
	blab["Content"] = template.HTML(string(content))
	blab["Metatitle"] = page["Metatitle"].(string)
	err = templates.ExecuteTemplate(w, "frame.html", &blab)

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
		// translate ContentUrl to Urlpath and redirect
		m, err := loadPages()
		if err != nil {
			return
		}
		// fill in content
		path := r.URL.Path
		pagename := Contenturl2pagename(path, m)
		page := m[pagename].(map[string]interface{})
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
