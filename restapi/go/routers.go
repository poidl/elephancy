/* OpenAPI spec version: 0.1.0
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
*
*modified by S. Riha (2017)
*/

package api

import "net/http"

type Cfh func(c Configuration, w http.ResponseWriter, r *http.Request)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc Cfh
}

type Routes []Route

var routes = Routes{
	Route{
		"ListPages",
		"GET",
		"/api/pages",
		ListPages,
	},
	// Route{
	// 	"FileServer",
	// 	"GET",
	// 	"/resources/",
	// 	FileServer,
	// },
	// Route{
	// Route{
	// 	"FindPageByPrettyURL",
	// 	"GET",
	// 	"/api/pages/FindPageByPrettyURL",
	// 	FindPageByPrettyURL,
	// },
	// 	"FindPageByLinksSelf",
	// 	"GET",
	// 	"/api/pages/FindPageByLinksSelf",
	// 	FindPageByLinksSelf,
	// },
	Route{
		"FindPageByKeyValue",
		"GET",
		"/api/pages/FindPageByKeyValue",
		FindPageByKeyValue,
	},
	Route{
		"GetPageContent",
		"GET",
		"/api/content/",
		GetPageContent,
	},
}

func MyRouter(c Configuration) http.Handler {
	sm := http.NewServeMux()
	for _, route := range routes {
		f := MyLogger(configureHandler(c, route.HandlerFunc), route.Name)
		sm.HandleFunc(route.Pattern, f)
	}
	return sm
}

func configureHandler(c Configuration, handlerfunc Cfh) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerfunc(c, w, r)
	}
}

// type Route struct {
// 	Name        string
// 	Method      string
// 	Pattern     string
// 	HandlerFunc http.HandlerFunc
// }

// type Routes []Route

// func NewRouter() *mux.Router {
// 	router := mux.NewRouter().StrictSlash(true)
// 	for _, route := range routes {
// 		var handler http.Handler
// 		handler = route.HandlerFunc
// 		handler = Logger(handler, route.Name)

// 		router.
// 			Methods(route.Method).
// 			Path(route.Pattern).
// 			Name(route.Name).
// 			Handler(handler)
// 	}

// 	return router
// }

// func Index(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello World!")
// }

// var routes = Routes{
// 	Route{
// 		"Index",
// 		"GET",
// 		"/api/",
// 		Index,
// 	},

// 	// Route{
// 	// 	"FindPageById",
// 	// 	"GET",
// 	// 	"/api/pages/{id}",
// 	// 	FindPageById,
// 	// },

// 	Route{
// 		"ListPages",
// 		"GET",
// 		"/api/pages",
// 		ListPages,
// 	},
// 	Route{
// 		"FindPageByPrettyURL",
// 		"GET",
// 		"/api/pages/FindPageByPrettyURL",
// 		FindPageByPrettyURL,
// 	},
// }
