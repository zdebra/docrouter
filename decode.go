package docrouter

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func DecodeQueryParams(structPtr interface{}, req *http.Request) error {
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
		queryParamName := tField.getTagName()
		if queryParamName == "" {
			continue
		}

		valueFromQuery := req.URL.Query().Get(queryParamName)
		structField := sElem.FieldByName(tField.name)
		if !structField.IsValid() {
			return fmt.Errorf("invalid field name %q", tField.name)
		}

		if !structField.CanSet() {
			return fmt.Errorf("can't set field %q", tField.name)
		}

		switch structField.Kind() {
		case reflect.Int:
			intVal, err := strconv.Atoi(valueFromQuery)
			if err != nil {
				return fmt.Errorf("converting %q to int: %v", valueFromQuery, err)
			}
			structField.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(valueFromQuery)
			if err != nil {
				return fmt.Errorf("converting %q to bool: %v", valueFromQuery, err)
			}
			structField.SetBool(boolVal)
		}

	}

	return nil
}
