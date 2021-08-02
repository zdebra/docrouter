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
	params := openapi3.Parameters{}
	if r.PathParams != nil {
		pathParams, err := createParamsWithReflection(r.PathParams, openapi3.ParameterInPath, true)
		if err != nil {
			return nil, fmt.Errorf("create path params: %w", err)
		}

		for _, pathParam := range pathParams {
			params = append(params, &openapi3.ParameterRef{
				Value: pathParam,
			})
		}
	}
	if r.QueryParams != nil {
		reflectedParams, err := createParamsWithReflection(r.QueryParams, openapi3.ParameterInQuery, false)
		if err != nil {
			return nil, fmt.Errorf("create query params: %w", err)
		}

		for _, param := range reflectedParams {
			params = append(params, &openapi3.ParameterRef{
				Value: param,
			})
		}
	}
	return params, nil
}

func createParamsWithReflection(structPtr interface{}, inParam string, forceRequired bool) ([]*openapi3.Parameter, error) {
	pParam, err := parseParameter(structPtr)
	if err != nil {
		return nil, fmt.Errorf("parsing param: %w", err)
	}

	params := []*openapi3.Parameter{}
	for _, tField := range pParam.fields {
		fieldName := tField.name
		var exampleTag interface{}
		var schemaType string
		switch tField.kind {
		case reflect.Int:
			x, err := strconv.Atoi(tField.getTagExample())
			if err != nil {
				return nil, fmt.Errorf("invalid int value for field %q, tag: `example`: %v", fieldName, err)
			}
			exampleTag = x
			schemaType = "integer"
		case reflect.Bool:
			x, err := strconv.ParseBool(tField.getTagExample())
			if err != nil {
				return nil, fmt.Errorf("invalid bool value for field %q, tag: `example`: %v", fieldName, err)
			}
			exampleTag = x
			schemaType = "boolean"
		default:
			exampleTag = tField.getTagExample()
			schemaType = "string"
		}

		required := true
		if !forceRequired {
			x, err := strconv.ParseBool(tField.getTagRequired())
			if err != nil {
				return nil, fmt.Errorf("invalid bool value for field %q, tag: `required`: %v", fieldName, err)
			}
			required = x
		}

		schemaFromTag, err := schemaFromTag(tField.getTagSchemaMin(), schemaType)
		if err != nil {
			return nil, fmt.Errorf("schemaFromTag: %w", err)
		}

		params = append(params, &openapi3.Parameter{
			Name:        tField.getTagName(),
			Description: tField.getTagDesc(),
			Example:     exampleTag,
			In:          inParam,
			Required:    required,
			Schema:      schemaFromTag,
		})
	}
	return params, nil
}

// todo expand this logic to accept more schemas from tags
func schemaFromTag(schemaMinTagValue string, schemaFieldType string) (*openapi3.SchemaRef, error) {
	if schemaMinTagValue == "" {
		return openapi3.NewSchemaRef("", &openapi3.Schema{
			Type: schemaFieldType,
		}), nil
	}
	fValue, err := strconv.ParseFloat(schemaMinTagValue, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing float from %q: %v", schemaMinTagValue, err)
	}
	return openapi3.NewSchemaRef("", &openapi3.Schema{
		Min:  &fValue,
		Type: schemaFieldType,
	}), nil
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
