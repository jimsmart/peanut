package peanut

import (
	"reflect"
)

type writer struct {
	records       []interface{}
	types         []reflect.Type
	headers       [][]string
	headersByType map[reflect.Type][]string
}

func (w *writer) init(x interface{}) bool {
	// Lazy init.
	if w.headersByType == nil {
		w.headersByType = make(map[reflect.Type][]string)
	}

	t := baseType(x)
	// Is this type already registered?
	if _, ok := w.headersByType[t]; ok {
		// Yes.
		return false
	}

	w.records = append(w.records, x)
	w.types = append(w.types, t)
	h := headers(x)
	w.headers = append(w.headers, h)
	w.headersByType[t] = h
	return true
}

func headers(x interface{}) []string {
	var out []string
	reflectStructFields(x, func(name string, t reflect.Type, tag string) {
		tag = firstTagValue(tag)
		out = append(out, tag)
	})
	return out
}
