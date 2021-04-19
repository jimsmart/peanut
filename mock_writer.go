package peanut

import (
	"fmt"
	"reflect"
)

var _ Writer = &MockWriter{}

// MockWriter captures written data in memory, to provide easy mocking
// when testing code that uses peanut.
type MockWriter struct {
	*base
	Headers            map[string][]string
	Data               map[string][]map[string]string
	DisableDataCapture map[string]bool
	CalledWrite        int
	CalledClose        int
	CalledCancel       int
}

func (w *MockWriter) register(x interface{}) (reflect.Type, error) {
	// Lazy init.
	if w.base == nil {
		w.base = &base{}
		w.Headers = make(map[string][]string)
		w.Data = make(map[string][]map[string]string)
	}
	// Register with base writer.
	t, ok := w.base.register(x)
	if !ok {
		return t, nil
	}
	w.Headers[t.Name()] = w.headersByType[t]
	return t, nil
}

// Write is called to persist records.
func (w *MockWriter) Write(x interface{}) error {
	w.CalledWrite++
	t, err := w.register(x)
	if err != nil {
		return err
	}

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
		out[tag] = fmt.Sprintf("%v", v)
	})
	return out
}
