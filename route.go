package docserver

import "net/http"

type Route struct {
	Path         string
	Methods      []string
	RequestBody  interface{}
	ResponseBody interface{}
	QueryParams  interface{}
	HeaderParams interface{}
	PathParams   interface{}
	Middlewares  []func(http.Handler) http.Handler
	Handler      http.Handler

	// Short summary
	Summary string
	// Optional description. Should use CommonMark syntax
	Description string
}
