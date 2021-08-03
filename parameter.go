package docrouter

import (
	"fmt"
	"reflect"
)

type parsedParameter struct {
	fields []taggedField
}

type taggedField struct {
	name               string
	kind               reflect.Kind
	rawTag             string
	parsedDocrouterTag map[string]string // "desc": "xxxx", "example": "3"
}

func parseParameter(structPtr interface{}) (parsedParameter, error) {
	var pParam parsedParameter
	v := reflect.ValueOf(structPtr).Elem()
	if !v.CanAddr() {
		return pParam, fmt.Errorf("item must be a pointer")
	}

	pParam.fields = []taggedField{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		docrouterTag, found := typeField.Tag.Lookup("docrouter")
		if !found {
			// field doesn't have a docrouter tag
			continue
		}
		fieldName := typeField.Name
		keys := []string{
			"kind",
			"name",
			"desc",
			"example",
			"required",
			"schemaMin",
		}
		parsedDocrouterTag := map[string]string{}
		for _, key := range keys {
			parsedDocrouterTag[key], _ = tagLookup(key, docrouterTag)
		}
		pParam.fields = append(pParam.fields, taggedField{
			name:               fieldName,
			kind:               typeField.Type.Kind(),
			rawTag:             docrouterTag,
			parsedDocrouterTag: parsedDocrouterTag,
		})

	}
	return pParam, nil
}

func (tf *taggedField) getTagExample() string {
	return tf.parsedDocrouterTag["example"]
}

func (tf *taggedField) getTagRequired() string {
	return tf.parsedDocrouterTag["required"]
}

func (tf *taggedField) getTagSchemaMin() string {
	return tf.parsedDocrouterTag["schemaMin"]
}

func (tf *taggedField) getTagName() string {
	return tf.parsedDocrouterTag["name"]
}

func (tf *taggedField) getTagDesc() string {
	return tf.parsedDocrouterTag["desc"]
}

func (tf *taggedField) getTagKind() string {
	return tf.parsedDocrouterTag["kind"]
}
