package myjson

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

type msi map[string]interface{}
type ia []Page
type page map[string]interface{}
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// LoadJSONmsi opens a json file and returns the contents as a map[string]interface{}
// TODO: handle errors
func LoadJSONmsi(filename string) (msi, error) {

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

func LoadJSONnew(filename string) (ia, error) {

	bytearr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var p []Page
	err = json.Unmarshal(bytearr, &p)
	if err != nil {
		return nil, err
	}
	// m := p.([]interface{})
	return p, nil
	// println(p[""])
	// return m, nil
}

func writeJson(filename string, msi map[string]interface{}) {

	data, err := json.Marshal(msi)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal("Writing " + filename + " failed.")
	}
}

// map2array orders page data into an array according linkweights. It takes a map[string]interface{}, which maps pagenames to pages. A page is a json object containing page data as unordered list of key:value pairs. It returns an array of these objects, which is sorted corresponding to the "Linkweight" keys of the page objects.
func (pcoll *msi) toArray() (arr []map[string]interface{}) {
	// we need an array for iteration order. see http://blog.golang.org/go-maps-in-action
	g := make(map[int]string)
	for pname, page := range *pcoll {
		pag := page.(map[string]interface{})
		g[int(pag["Linkweight"].(float64))] = pname
	}

	var keys []int
	for k := range g {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		arr = append(arr, (*pcoll)[g[k]].(map[string]interface{}))
	}
	return arr
}

// TODO: urlPathToPagename and contentURLToPagename almost same

func (pcoll *msi) urlPathToPage(urlpath string) (page, error) {
	// TODO: return error in case it doesn't find anything
	for _, page := range *pcoll {
		pag := page.(map[string]interface{})
		if urlpath == pag["Urlpath"].(string) {
			return pag, nil
		}
	}
	return nil, &errorString{"Page not found"}
}

func (pcoll *msi) ContentURLToPage(contenturl string) (page, error) {
	// TODO: return error in case it doesn't find anything
	for _, page := range *pcoll {
		pag := page.(map[string]interface{})
		if contenturl == pag["ContentUrl"].(string) {
			return pag, nil
		}
	}
	return nil, &errorString{"Page not found"}
}

func (pcoll *ia) PrettyURLToPage(prettyURL string) (Page, error) {
	// TODO: return error in case it doesn't find anything
	for _, page := range *pcoll {
		if "/"+prettyURL == page.Prettyurl {
			return page, nil
		}
	}
	return Page{}, &errorString{"Page not found"}
}

func ContentURLToUrlpath(contenturl string) (string, error) {
	filename := "./json/pages.json"
	pcoll, err := LoadJSONmsi(filename)
	if err != nil {
		return "", err
	}
	page, err := pcoll.ContentURLToPage(contenturl)
	urlpath := page["Urlpath"].(string)
	return urlpath, nil
}

func (pcoll *msi) getPage(pagename string) (page, error) {
	// load the page data into a map[string]interface{}
	pg := (*pcoll)[pagename].(map[string]interface{})
	if pg == nil {
		return nil, &errorString{"Page not found"}
	}
	return pg, nil
}

func (pcoll *msi) contentFromPage(pg page) (content []byte, modtime time.Time, err error) {
	// read content from html file and get modification time
	content, modtime, err = readContent(pg["ContentUrl"].(string))
	return content, modtime, err
}

// func getTemplateData(urlpath string) (map[string]interface{}, time.Time, error) {
// 	filename := "./json/pages.json"
// 	pcoll, err := LoadJSONmsi(filename)
// 	if err != nil {
// 		return nil, time.Time{}, err
// 	}
// 	page, err := pcoll.urlPathToPage(urlpath)
// 	if err != nil {
// 		return nil, time.Time{}, err
// 	}
// 	content, modtime, err := pcoll.contentFromPage(page)
// 	arr := pcoll.toArray()
// 	blab := make(map[string]interface{})
// 	blab["Pages"] = arr
// 	blab["Content"] = template.HTML(string(content))
// 	blab["Metatitle"] = page["Metatitle"].(string)
// 	return blab, modtime, err
// }

func GetTemplateData(urlpath string) (map[string]interface{}, time.Time, error) {
	filename := "./json/pages.json"
	pcoll, err := LoadJSONmsi(filename)
	if err != nil {
		return nil, time.Time{}, err
	}
	page, err := pcoll.urlPathToPage(urlpath)
	if err != nil {
		return nil, time.Time{}, err
	}
	content, modtime, err := pcoll.contentFromPage(page)
	arr := pcoll.toArray()
	blab := make(map[string]interface{})
	blab["Pages"] = arr
	blab["Content"] = template.HTML(string(content))
	blab["Metatitle"] = page["Metatitle"].(string)
	return blab, modtime, err
}

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