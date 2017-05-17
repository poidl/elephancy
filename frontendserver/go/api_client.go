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

package frontendserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func callAPI(path string,
	method string,
	headerParams map[string]string,
	queryParams url.Values) (*http.Response, error) {
	switch strings.ToUpper(method) {
	case "GET":
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			log.Fatal(err)
		}
		// add header parameter, if any
		if len(headerParams) > 0 {
			for key, value := range headerParams {
				req.Header.Add(key, value)
			}
		}
		req.URL.RawQuery = queryParams.Encode()
		// println(req.URL.String())
		response, err := http.DefaultClient.Do(req)

		return response, err
		// case "POST":
		// 	response, err := request.Post(path)
		// 	return response, err
		// case "PUT":
		// 	response, err := request.Put(path)
		// 	return response, err
		// case "PATCH":
		// 	response, err := request.Patch(path)
		// 	return response, err
		// case "DELETE":
		// 	response, err := request.Delete(path)
		// 	return response, err
	}

	return nil, fmt.Errorf("invalid method %v", method)
}
