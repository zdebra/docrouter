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
		paramName, paramKind := tField.getTagName(), tField.getTagKind()
		if paramName == "" {
			continue
		}

		valueStr, err := strValueFromRequest(paramName, paramKind, req)
		if err != nil {
			return fmt.Errorf("read string value from request: %v", err)
		}

		if err := convertAndSetStructField(&sElem, tField.name, valueStr); err != nil {
			return fmt.Errorf("convert and set struct field: %w", err)
		}
	}

	return nil
}

func strValueFromRequest(paramName, kind string, req *http.Request) (string, error) {
	switch kind {
	case openapi3.ParameterInQuery:
		return req.URL.Query().Get(paramName), nil
	case openapi3.ParameterInPath:
		return mux.Vars(req)[paramName], nil
	case openapi3.ParameterInCookie:
		c, err := req.Cookie(paramName)
		if err != nil {
			return "", nil
		}
		return c.Value, nil
	case openapi3.ParameterInHeader:
		return req.Header.Get(paramName), nil
	default:
		return "", fmt.Errorf("paramter kind %q not supported", kind)
	}
}

func convertAndSetStructField(sVal *reflect.Value, fieldName, valueStr string) error {
	structField := sVal.FieldByName(fieldName)
	if !structField.IsValid() {
		return fmt.Errorf("invalid field name %q", fieldName)
	}

	if !structField.CanSet() {
		return fmt.Errorf("can't set field %q", fieldName)
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
	return nil
}
