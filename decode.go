package docrouter

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

func DecodeParams(structPtr interface{}, req *http.Request) error {
	if req.URL == nil {
		return fmt.Errorf("invalid request - req.URL is nil")
	}

	pParam, err := parseParameter(structPtr)
	if err != nil {
		return fmt.Errorf("parsing param: %w", err)
	}

	ps := reflect.ValueOf(structPtr)
	sElem := ps.Elem()
	if sElem.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct pointer")
	}

	for _, tField := range pParam.fields {
		paramName := tField.getTagName()
		if paramName == "" {
			continue
		}

		valueStr := ""
		switch tField.getTagKind() {
		case openapi3.ParameterInQuery:
			valueStr = req.URL.Query().Get(paramName)
		case openapi3.ParameterInPath:
			x, found := mux.Vars(req)[paramName]
			if !found {
				return fmt.Errorf("parameter %q not found in path", paramName)
			}
			valueStr = x
		default:
			return fmt.Errorf("paramter kind %q not supported", tField.getTagKind())
		}

		structField := sElem.FieldByName(tField.name)
		if !structField.IsValid() {
			return fmt.Errorf("invalid field name %q", tField.name)
		}

		if !structField.CanSet() {
			return fmt.Errorf("can't set field %q", tField.name)
		}

		switch structField.Kind() {
		case reflect.Int:
			intVal, err := strconv.Atoi(valueStr)
			if err != nil {
				return fmt.Errorf("converting %q to int: %v", valueStr, err)
			}
			structField.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(valueStr)
			if err != nil {
				return fmt.Errorf("converting %q to bool: %v", valueStr, err)
			}
			structField.SetBool(boolVal)
		case reflect.String:
			structField.SetString(valueStr)
		default:
			return fmt.Errorf("unsupported conversion for %v", structField.Kind())
		}

	}

	return nil
}
