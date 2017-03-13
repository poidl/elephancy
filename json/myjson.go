package myjson

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type msi map[string]interface{}
type Pages []Page
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

func WriteJson(filename string, msi map[string]interface{}) {
	data, err := json.Marshal(msi)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal("Writing " + filename + " failed.")
	}
}

func LoadPages(filename string) (Pages, error) {

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

// // map2array orders page data into an array according linkweights. It takes a map[string]interface{}, which maps pagenames to pages. A page is a json object containing page data as unordered list of key:value pairs. It returns an array of these objects, which is sorted corresponding to the "Linkweight" keys of the page objects.
// func (pcoll *msi) toArray() (arr []map[string]interface{}) {
// 	// we need an array for iteration order. see http://blog.golang.org/go-maps-in-action
// 	g := make(map[int]string)
// 	for pname, page := range *pcoll {
// 		pag := page.(map[string]interface{})
// 		g[int(pag["Linkweight"].(float64))] = pname
// 	}

// 	var keys []int
// 	for k := range g {
// 		keys = append(keys, k)
// 	}
// 	sort.Ints(keys)

// 	for _, k := range keys {
// 		arr = append(arr, (*pcoll)[g[k]].(map[string]interface{}))
// 	}
// 	return arr
// }

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

// func (pcoll *msi) ContentURLToPage(contenturl string) (page, error) {
// 	// TODO: return error in case it doesn't find anything
// 	for _, page := range *pcoll {
// 		pag := page.(map[string]interface{})
// 		if contenturl == pag["ContentUrl"].(string) {
// 			return pag, nil
// 		}
// 	}
// 	return nil, &errorString{"Page not found"}
// }

func (pages *Pages) GetPageByPrettyURL(prettyURL string) (Page, error) {
	for _, page := range *pages {
		if prettyURL == page.Prettyurl {
			return page, nil
		}
	}
	return Page{}, &errorString{"Page not found"}
}

func (pages *Pages) GetPageByLinksSelf(linksSelf string) (Page, error) {
	for _, page := range *pages {
		if linksSelf == page.Links.Self {
			return page, nil
		}
	}
	return Page{}, &errorString{"Page not found"}
}
