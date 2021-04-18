package peanut

import (
	"fmt"
	"reflect"
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
		case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// Allowed type: append stringified value to list.
			out = append(out, fmt.Sprintf("%v", v))
		default:
			m := fmt.Sprintf("Unknown type: %v", v) // TODO(js) This would be clearer if it used t.Name() ?
			panic(m)
		}
	})
	return out
}

func mapValues(x interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	// TODO(js) Can we refactor reflect-type->value out, to reduce cyclomatic complexity?
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		tag = firstTagValue(tag)
		switch t.Kind() {
		case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// Allowed type: put value into map.
			out[tag] = v
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
