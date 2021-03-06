package peanut

import (
	"reflect"
)

type base struct {
	headersByType map[reflect.Type][]string       // headersByType is a list of headers for each struct type.
	typesByType   map[reflect.Type][]reflect.Type // typesByType is a list of reflected field types for each struct type.
	tagsByType    map[reflect.Type][]string       // tagsByType is a list of field tags for each struct type.
	closed        bool
}

// register a type and collect its metadata.
// If the type is a newly registered type
// (has not been seen before),
// return true. Otherwise return false.
func (w *base) register(x interface{}) (reflect.Type, bool) {
	// Lazy init.
	if w.headersByType == nil {
		w.headersByType = make(map[reflect.Type][]string)
		w.typesByType = make(map[reflect.Type][]reflect.Type)
		w.tagsByType = make(map[reflect.Type][]string)
	}

	t := baseType(x)
	// Is this type already registered?
	if _, ok := w.headersByType[t]; ok {
		// Yes.
		return t, false
	}

	var headers []string
	var types []reflect.Type
	var tags []string

	reflectStructFields(x, func(name string, t reflect.Type, tag string) {
		h := firstTagValue(tag)
		headers = append(headers, h)
		types = append(types, t)
		tags = append(tags, tag)
	})

	w.headersByType[t] = headers
	w.typesByType[t] = types
	w.tagsByType[t] = tags
	return t, true
}
