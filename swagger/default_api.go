/*
 * Elephancy backend (Simple)
 *
 * Elephancy backend (Simple)
 *
 * OpenAPI spec version: 0.1.0
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
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

package swagger

import (
	"bytes"
	"encoding/json"
	"log"
)

// type DefaultApi struct {
// 	Configuration Configuration
// }

// func NewDefaultApi() *DefaultApi {
// 	configuration := NewConfiguration()
// 	return &DefaultApi{
// 		Configuration: *configuration,
// 	}
// }

// func NewDefaultApiWithBasePath(basePath string) *DefaultApi {
// 	configuration := NewConfiguration()
// 	configuration.BasePath = basePath

// 	return &DefaultApi{
// 		Configuration: *configuration,
// 	}
// }

// /**
//  *
//  * Returns a page based on a single ID
//  *
//  * @param id ID of page to fetch
//  * @return *Page
//  */
// func (a DefaultApi) FindPageById(id int64) (*Page, *APIResponse, error) {

// 	var httpMethod = "Get"
// 	// create path and map variables
// 	path := a.Configuration.BasePath + "/pages/{id}"
// 	path = strings.Replace(path, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

// 	headerParams := make(map[string]string)
// 	queryParams := url.Values{}
// 	formParams := make(map[string]string)
// 	var postBody interface{}
// 	var fileName string
// 	var fileBytes []byte
// 	// add default headers if any
// 	for key := range a.Configuration.DefaultHeader {
// 		headerParams[key] = a.Configuration.DefaultHeader[key]
// 	}

// 	// to determine the Content-Type header
// 	localVarHTTPContentTypes := []string{"application/json"}

// 	// set Content-Type header
// 	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHTTPContentTypes)
// 	if localVarHttpContentType != "" {
// 		headerParams["Content-Type"] = localVarHttpContentType
// 	}
// 	// to determine the Accept header
// 	localVarHttpHeaderAccepts := []string{
// 		"application/json",
// 		"application/xml",
// 		"text/xml",
// 		"text/html",
// 	}

// 	// set Accept header
// 	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
// 	if localVarHttpHeaderAccept != "" {
// 		headerParams["Accept"] = localVarHttpHeaderAccept
// 	}
// 	var successPayload = new(Page)
// 	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
// 	if err != nil {
// 		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
// 	}
// 	err = json.Unmarshal(httpResponse.Body(), &successPayload)
// 	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
// }

/**
 *
 * Returns all pages
 *
 * @return []Page
 */

var basepath = "http://127.0.0.1:8088/api"

func ListPages() ([]Page, error) {

	var httpMethod = "GET"
	// create path and map variables
	path := basepath + "/pages"

	var successPayload = new([]Page)
	httpResponse, err := callAPI(path, httpMethod)
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
