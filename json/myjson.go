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
