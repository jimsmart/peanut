package peanut

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Writer defines a record-based writer.
type Writer interface {
	Write(r interface{}) error
	Close() error
	Cancel() error
}

//

const tagName = "peanut"

func stringValues(x interface{}) []string {
	var out []string
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		switch t.Kind() {
		case reflect.String:
			out = append(out, v.(string))
		case reflect.Int:
			out = append(out, strconv.FormatInt(int64(v.(int)), 10))
		case reflect.Bool:
			out = append(out, strconv.FormatBool(v.(bool)))
		case reflect.Float64, reflect.Float32:
			out = append(out, fmt.Sprint(v.(float64))) // TODO(js) This seems lazy.
		case reflect.Int8:
			out = append(out, strconv.FormatInt(int64(v.(int8)), 10))
		case reflect.Int16:
			out = append(out, strconv.FormatInt(int64(v.(int16)), 10))
		case reflect.Int32:
			out = append(out, strconv.FormatInt(int64(v.(int32)), 10))
		case reflect.Int64:
			out = append(out, strconv.FormatInt(v.(int64), 10))
		case reflect.Uint:
			out = append(out, strconv.FormatUint(uint64(v.(uint)), 10))
		case reflect.Uint8:
			out = append(out, strconv.FormatUint(uint64(v.(uint8)), 10))
		case reflect.Uint16:
			out = append(out, strconv.FormatUint(uint64(v.(uint16)), 10))
		case reflect.Uint32:
			out = append(out, strconv.FormatUint(uint64(v.(uint32)), 10))
		case reflect.Uint64:
			out = append(out, strconv.FormatUint(v.(uint64), 10))
		default:
			m := fmt.Sprintf("Unknown type: %v", v) // TODO(js) This would be clearer if it used t.Name() ?
			panic(m)
		}
	})
	return out
}

func mapValues(x interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		tag = firstTagValue(tag)
		switch t.Kind() {
		case reflect.String:
			out[tag] = v.(string)
		case reflect.Int:
			out[tag] = v.(int)
		case reflect.Bool:
			out[tag] = v.(bool)
		case reflect.Float64:
			out[tag] = v.(float64)
		case reflect.Float32:
			out[tag] = v.(float32)
		case reflect.Int8:
			out[tag] = v.(int8)
		case reflect.Int16:
			out[tag] = v.(int16)
		case reflect.Int32:
			out[tag] = v.(int32)
		case reflect.Int64:
			out[tag] = v.(int64)
		case reflect.Uint8:
			out[tag] = v.(uint8)
		case reflect.Uint16:
			out[tag] = v.(uint16)
		case reflect.Uint32:
			out[tag] = v.(uint32)
		case reflect.Uint64:
			out[tag] = v.(uint64)
		case reflect.Uint:
			out[tag] = v.(uint)
		default:
			m := fmt.Sprintf("Unknown type: %v", v)
			panic(m)
		}
	})
	return out
}

func reflectStructFields(x interface{}, fn func(name string, t reflect.Type, tag string)) {

	// TODO(js) This should work with Ptr and non-Ptr.

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := baseType(x)

	// Get the type and kind of our user variable.
	// fmt.Println("Type:", t.Name())
	// fmt.Println("Kind:", t.Kind())

	// Iterate over all available fields and read the tag value.
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)
		tag := field.Tag.Get(tagName)

		// Only process fields with appropriate tags.
		if tag != "" {
			fn(field.Name, field.Type, tag)
		}

		// fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}

func baseType(x interface{}) reflect.Type {
	// TODO(js) This should work with Ptr and non-Ptr.
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func reflectStructValues(x interface{}, fn func(name string, t reflect.Type, v interface{}, tag string)) {

	// TODO(js) This should work with Ptr and non-Ptr.

	t := reflect.ValueOf(x)
	// TODO This was taken from elsewhere, and modified - but can t.Kind even be a ptr here?
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)
		tag := field.Tag.Get(tagName)

		// Only process fields with appropriate tags.
		if tag != "" {

			name := field.Name
			// Filter out unexported fields.
			r, _ := utf8.DecodeRuneInString(name)
			if !unicode.IsUpper(r) {
				continue
			}
			val := t.Field(i).Interface()

			fn(name, field.Type, val, tag)
		}

		// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", field.Name, val, tag)
	}
}

func firstTagValue(s string) string {
	return strings.Split(s, ",")[0]
}

func secondTagValue(s string) string {
	x := strings.Split(s, ",")
	if len(x) > 1 {
		return x[1]
	}
	return ""
}
