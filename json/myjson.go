package myjson

import (
	"encoding/json"
	"io/ioutil"
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

func (pages *Pages) GetPageById(id int64) (Page, error) {
	for _, page := range *pages {
		if id == page.Id {
			return page, nil
		}
	}
	return Page{}, &errorString{"Page not found"}
}

// func (pages *Pages) GetPageByLinksSelf(linksSelf string) (Page, error) {
// 	for _, page := range *pages {
// 		link, _ := page.GetLinkByRel("self")
// 		if linksSelf == link {
// 			return page, nil
// 		}
// 	}
// 	return Page{}, &errorString{"Page not found"}
// }

func (page *Page) GetLinkByRel(rel string) (link string, e error) {
	for _, link := range (*page).Links {
		if link.Rel == rel {
			return link.Href, nil
		}
	}
	return "", &errorString{"Link not found"}
}
