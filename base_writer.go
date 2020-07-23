package peanut

import (
	"reflect"
)

type base struct {
	records       []interface{}
	types         []reflect.Type
	typeHeaders   [][]string
	typeTypes     [][]reflect.Type
	headersByType map[reflect.Type][]string
	typesByType   map[reflect.Type][]reflect.Type
}

// register a type and collect its metadata.
// If the type is a newly registered type
// (has not been seen before),
// return true. Otherwise return false.
func (w *base) register(x interface{}) bool {
	// Lazy init.
	if w.headersByType == nil {
		w.headersByType = make(map[reflect.Type][]string)
		w.typesByType = make(map[reflect.Type][]reflect.Type)
	}

	t := baseType(x)
	// Is this type already registered?
	if _, ok := w.headersByType[t]; ok {
		// Yes.
		return false
	}

	w.records = append(w.records, x)
	w.types = append(w.types, t)

	var headers []string
	var types []reflect.Type

	reflectStructFields(x, func(name string, t reflect.Type, tag string) {
		tag = firstTagValue(tag)
		headers = append(headers, tag)
		types = append(types, t)
	})

	w.typeHeaders = append(w.typeHeaders, headers)
	w.headersByType[t] = headers
	w.typeTypes = append(w.typeTypes, types)
	w.typesByType[t] = types
	return true
}
