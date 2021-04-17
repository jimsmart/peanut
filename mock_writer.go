package peanut

import (
	"fmt"
	"reflect"
	"strconv"
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
		switch t.Kind() {
		case reflect.String:
			out[tag] = v.(string)
		case reflect.Int:
			out[tag] = strconv.FormatInt(int64(v.(int)), 10)
		case reflect.Bool:
			out[tag] = strconv.FormatBool(v.(bool))
		case reflect.Float64, reflect.Float32:
			out[tag] = fmt.Sprint(v.(float64)) // TODO(js) This seems lazy.
		case reflect.Int8:
			out[tag] = strconv.FormatInt(int64(v.(int8)), 10)
		case reflect.Int16:
			out[tag] = strconv.FormatInt(int64(v.(int16)), 10)
		case reflect.Int32:
			out[tag] = strconv.FormatInt(int64(v.(int32)), 10)
		case reflect.Int64:
			out[tag] = strconv.FormatInt(v.(int64), 10)
		case reflect.Uint:
			out[tag] = strconv.FormatUint(uint64(v.(uint)), 10)
		case reflect.Uint8:
			out[tag] = strconv.FormatUint(uint64(v.(uint8)), 10)
		case reflect.Uint16:
			out[tag] = strconv.FormatUint(uint64(v.(uint16)), 10)
		case reflect.Uint32:
			out[tag] = strconv.FormatUint(uint64(v.(uint32)), 10)
		case reflect.Uint64:
			out[tag] = strconv.FormatUint(v.(uint64), 10)
		default:
			m := fmt.Sprintf("Unknown type: %v", v)
			panic(m)
		}
	})
	return out
}
