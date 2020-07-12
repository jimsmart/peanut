package peanut

import (
	"fmt"
	"reflect"
	"strconv"
)

var _ Writer = &MockWriter{}

// MockWriter is a type of peanut.Writer that captures data, for testing purposes.
type MockWriter struct {
	*writer
	Headers            map[string][]string
	Data               map[string][]map[string]string
	DisableDataCapture map[string]bool
	CalledWrite        int
	CalledClose        int
	CalledCancel       int
}

func (w *MockWriter) initialise(x interface{}) error {
	// Lazy init.
	if w.writer == nil {
		w.writer = &writer{}
		w.Headers = make(map[string][]string)
		w.Data = make(map[string][]map[string]string)
	}
	// Register with base writer.
	ok := w.init(x)
	if !ok {
		return nil
	}
	i := len(w.types) - 1
	w.Headers[w.types[i].Name()] = w.headers[i]
	return nil
}

// Write is called to persist records.
func (w *MockWriter) Write(x interface{}) error {
	w.CalledWrite++
	err := w.initialise(x)
	if err != nil {
		return err
	}
	t := baseType(x)

	n := t.Name()
	if w.DisableDataCapture != nil && w.DisableDataCapture[n] {
		return nil
	}
	// Capture the row of data.
	data := w.Data[n]
	data = append(data, stringValuesAsMap(x))
	w.Data[n] = data
	return nil
}

// Close should be called after successfully writing records.
func (w *MockWriter) Close() error {
	w.CalledClose++
	return nil
}

// Cancel should be called in the event of an error occurring.
func (w *MockWriter) Cancel() error {
	w.CalledCancel++
	return nil
}

func stringValuesAsMap(x interface{}) map[string]string {
	out := make(map[string]string)
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		tag = firstTagValue(tag)
		switch t.Kind() {
		case reflect.String:
			// out = append(out, v.(string))
			out[tag] = v.(string)
		case reflect.Int:
			// out = append(out, strconv.Itoa(v.(int)))
			out[tag] = strconv.Itoa(v.(int))
		default:
			m := fmt.Sprintf("Unknown type: %v", v)
			panic(m)
		}
	})
	return out
}
