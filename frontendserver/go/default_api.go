/*
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Modified by S. Riha
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package frontendserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	jc "github.com/poidl/elephancy/jsoncommon"
)

/**
 *
 * Returns all pages
 *
 * @return []Page
 */

var basepath = "http://127.0.0.1:8088/api"

// ListPages lists pages
func ListPages() (jc.Pages, error) {
	var httpMethod = "GET"
	// create path and map variables
	path := basepath + "/pages"
	var successPayload = new([]jc.Page)
	headerParams := make(map[string]string)
	httpResponse, err := callAPI(path, httpMethod, headerParams, url.Values{})
	if err != nil {
		return *successPayload, err
	}
	defer httpResponse.Body.Close()
	var b bytes.Buffer
	_, err = b.ReadFrom(httpResponse.Body)
	if err != nil {
		log.Fatal(err)

	}
	err = json.Unmarshal(b.Bytes(), &successPayload)
	return *successPayload, err
}

// // FindPageByPrettyURL returns page ID
// func FindPageByPrettyURL(prettyurl string) (jc.Page, error) {
// 	var httpMethod = "GET"
// 	// create path and map variables
// 	path := basepath + "/pages/FindPageByPrettyURL"

// 	var successPayload = new(jc.Page)
// 	queryParams := url.Values{}
// 	queryParams.Set("prettyurl", prettyurl)
// 	httpResponse, err := callAPI(path, httpMethod, queryParams)
// 	if err != nil {
// 		return *successPayload, err
// 	}
// 	defer httpResponse.Body.Close()
// 	if httpResponse.StatusCode == 404 {
// 		return *successPayload, fmt.Errorf("Page not found")
// 	}
// 	var b bytes.Buffer
// 	_, err = b.ReadFrom(httpResponse.Body)
// 	if err != nil {
// 		log.Fatal(err)

// 	}
// 	err = json.Unmarshal(b.Bytes(), &successPayload)
// 	return *successPayload, err
// }

// func FindPageByLinksSelf(link string) (jc.Page, error) {
// 	var httpMethod = "GET"
// 	// create path and map variables
// 	path := basepath + "/pages/FindPageByLinksSelf"
// 	println("****BELLO")
// 	var successPayload = new(jc.Page)
// 	queryParams := url.Values{}
// 	queryParams.Set("prettyurl", link)
// 	httpResponse, err := callAPI(path, httpMethod, queryParams)
// 	if err != nil {
// 		return *successPayload, err
// 	}
// 	defer httpResponse.Body.Close()
// 	if httpResponse.StatusCode == 404 {
// 		return *successPayload, fmt.Errorf("Page not found")
// 	}
// 	var b bytes.Buffer
// 	_, err = b.ReadFrom(httpResponse.Body)
// 	if err != nil {
// 		log.Fatal(err)

// 	}
// 	err = json.Unmarshal(b.Bytes(), &successPayload)
// 	return *successPayload, err
// }

func FindPageByKeyValue(key string, value string) (jc.Page, error) {
	var httpMethod = "GET"
	// create path and map variables
	path := basepath + "/pages/FindPageByKeyValue"
	var successPayload = new(jc.Page)
	queryParams := url.Values{}
	queryParams.Set("key", key)
	queryParams.Add("value", value)
	headerParams := make(map[string]string)
	httpResponse, err := callAPI(path, httpMethod, headerParams, queryParams)
	if err != nil {
		return *successPayload, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == 404 {
		return *successPayload, fmt.Errorf("Page not found")
	}
	var b bytes.Buffer
	_, err = b.ReadFrom(httpResponse.Body)
	if err != nil {
		log.Fatal(err)

	}
	err = json.Unmarshal(b.Bytes(), &successPayload)
	return *successPayload, err
}

func GetPageContent(id int64) (string, time.Time, error) {
	var httpMethod = "GET"
	// create path and map variables
	path := basepath + "/content/" + strconv.FormatInt(id, 10)
	var successPayload = new(string)
	queryParams := url.Values{}
	headerParams := make(map[string]string)
	headerParams["myheader"] = "XMLHttpRequest"
	httpResponse, err := callAPI(path, httpMethod, headerParams, queryParams)
	if err != nil {
		return *successPayload, time.Time{}, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == 404 {
		return *successPayload, time.Time{}, fmt.Errorf("Page not found")
	}
	var b bytes.Buffer
	_, err = b.ReadFrom(httpResponse.Body)
	if err != nil {
		log.Fatal(err)

	}
	// check when content was last modified
	lastmodified, _ := http.ParseTime(httpResponse.Header.Get("Last-Modified"))
	return b.String(), lastmodified, nil
}
