package docrouter

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

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

func (r *Route) openAPI3Params() (openapi3.Parameters, error) {
	// Path: name,description,example
	params := openapi3.NewParameters()
	if r.PathParams != nil {
		v := reflect.ValueOf(r.PathParams).Elem()
		if !v.CanAddr() {
			return nil, fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
		}

		for i := 0; i < v.NumField(); i++ {
			typeField := v.Type().Field(i)
			docrouterTag, found := typeField.Tag.Lookup("docrouter")
			if !found {
				// field doesn't have a docrouter tag
				continue
			}
			// fieldName := typeField.Name
			nameTag, _ := tagLookup("name", docrouterTag)
			descTag, _ := tagLookup("desc", docrouterTag)
			exampleStrTag, _ := tagLookup("example", docrouterTag)

			var exampleTag interface{}
			switch typeField.Type.Kind() {
			case reflect.Int:
				exampleTag, _ = strconv.Atoi(exampleStrTag)
			default:
				exampleTag = exampleStrTag

			}

			pathParam := openapi3.NewPathParameter(nameTag)
			pathParam.Description = descTag
			pathParam.Example = exampleTag

			params = append(params, &openapi3.ParameterRef{
				Value: pathParam,
			})
		}
	}
	return params, nil
}

func tagLookup(fieldName, rawTag string) (string, bool) {
	fieldName = fieldName + ":"
	if !strings.Contains(rawTag, fieldName) {
		return "", false
	}
	splits := strings.Split(rawTag, fieldName)
	if len(splits) < 2 {
		return "", false
	}
	valueWithTail := splits[1]

	val := valueWithTail
	if idx := strings.Index(valueWithTail, ";"); idx != -1 {
		val = valueWithTail[:idx]
	}
	return strings.TrimSpace(val), true
}
