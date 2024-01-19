package mysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

const (
	defaultDateFormatQuery = "%Y-%m-%dT%TZ"
	msgErrorRunTemplate    = "failed run template with model: %+v and template: %v"
)

func runTemplate(def string, i interface{}) (string, error) {
	if i == nil {
		return def, nil
	}

	funcMap := sprig.GenericFuncMap()
	funcMap["joins"] = Joins
	funcMap["defaultDateFormat"] = defaultDateFormat
	funcMap["searchLike"] = searchLike
	funcMap["parseJson"] = parseJson
	funcMap["convertTz"] = convertTz

	t := template.Must(template.New("").
		Funcs(funcMap).Parse(def))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, escapeSQLString(i)); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// created new interface with same type, but all string data is single quote escaped
func escapeSQLString(data interface{}) interface{} {
	if data == nil {
		return data
	}

	oldData := reflect.ValueOf(data)
	newData := reflect.New(oldData.Type()).Elem()
	escapeSingleQuote(newData, oldData)

	return newData.Interface()
}

func escapeSingleQuote(newData, oldData reflect.Value) {
	switch oldData.Type().Name() {
	case "string":
		newData.SetString(strings.Replace(oldData.String(), "'", "''", -1))
	case "LiteralString":
		newData.Set(oldData)
	case "Time":
		newData.Set(oldData)
	default:
		switch oldData.Kind() {
		case reflect.Ptr:
			pointedValue := oldData.Elem()
			if !pointedValue.IsValid() {
				return
			}
			newData.Set(reflect.New(pointedValue.Type()))
			escapeSingleQuote(newData.Elem(), pointedValue)
		case reflect.Interface:
			pointedValue := oldData.Elem()
			newValue := reflect.New(pointedValue.Type()).Elem()
			escapeSingleQuote(newValue, pointedValue)
			newData.Set(newValue)
		case reflect.Slice:
			newData.Set(reflect.MakeSlice(oldData.Type(), oldData.Len(), oldData.Cap()))
			for i := 0; i < oldData.Len(); i++ {
				escapeSingleQuote(newData.Index(i), oldData.Index(i))
			}
		case reflect.Map:
			newData.Set(reflect.MakeMap(oldData.Type()))
			for _, key := range oldData.MapKeys() {
				originalValue := oldData.MapIndex(key)
				newValue := reflect.New(originalValue.Type()).Elem()
				escapeSingleQuote(newValue, originalValue)
				newData.SetMapIndex(key, newValue)
			}
		case reflect.Struct:
			for i := 0; i < oldData.NumField(); i++ {
				escapeSingleQuote(newData.Field(i), oldData.Field(i))
			}
		default:
			newData.Set(oldData)
		}
	}
}

func Joins(in interface{}, prefix, suffix, delim string) (result string) {
	inType := reflect.TypeOf(in)
	inData := reflect.ValueOf(in)

	inKind := inType.Kind()
	if in == nil || (inKind != reflect.Slice && inKind != reflect.Array) {
		return result
	}

	for i := 0; i < inData.Len(); i++ {
		if i > 0 {
			result += delim
		}
		result += fmt.Sprintf("%v%v%v", prefix, inData.Index(i), suffix)
	}
	return
}

func defaultDateFormat(in string) string {
	return dateFormat(in, defaultDateFormatQuery)
}

func dateFormat(in, format string) string {
	return fmt.Sprintf("DATE_FORMAT(%v, '%v')", in, format)
}

func searchLike(in string) string {
	return fmt.Sprintf("'%v'", "%"+in+"%")
}

func parseJson(in interface{}) string {
	s, _ := json.Marshal(in)
	return string(s)
}

func convertTz(in string) string {
	return fmt.Sprintf("CONVERT_TZ('%v 00:00:00', '+07:00', '+00:00')", in)
}
