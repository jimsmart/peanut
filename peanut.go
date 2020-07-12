package peanut

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Writer defines a record-based writer.
type Writer interface {
	// Write is called to persist records.
	Write(r interface{}) error
	Close() error
	Cancel() error
}

//

const tagName = "peanut"

func reflectStructFields(x interface{}, fn func(name string, t reflect.Type, tag string)) {

	// TODO This should work with Ptr and non-Ptr.

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

		// Get the field tag value.
		tag := field.Tag.Get(tagName)

		// TODO Handle missing tag?

		// fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		fn(field.Name, field.Type, tag)
	}
}

func baseType(x interface{}) reflect.Type {
	// TODO This should work with Ptr and non-Ptr.
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func reflectStructValues(x interface{}, fn func(name string, t reflect.Type, v interface{}, tag string)) {

	// TODO This should work with Ptr and non-Ptr.

	t := reflect.ValueOf(x)
	// TODO This was taken from elsewhere, and modified - but can t.Kind even be a ptr here?
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)
		name := field.Name
		// Filter out unexported fields.
		r, _ := utf8.DecodeRuneInString(name)
		if !unicode.IsUpper(r) {
			continue
		}
		val := t.Field(i).Interface()
		tag := field.Tag.Get(tagName)

		// fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", field.Name, val, tag)
		fn(name, field.Type, val, tag)
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
