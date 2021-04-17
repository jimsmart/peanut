package peanut

import (
	"fmt"
	"reflect"
)

var _ Writer = &ExcelWriter{}

// ExcelWriter writes records to Excel files, writing
// each record type to an individual Excel file automatically.
//
// Filenames for each corresponding record type are derived
// accordingly:
//  prefix + type.Name() + suffix + ".xslx"
//
// The first row of resulting Excel file(s) will contain
// headers using names extracted from the struct's
// field tags, and will be frozen. Records' fields are
// written in the order that they appear within the struct.
//
// The caller must call Close on successful completion
// of all writing, to ensure buffers are flushed and
// files are properly written to disk.
//
// In the event of an error or cancellation, the
// caller must call Cancel before quiting, to ensure
// closure and cleanup of any partially written files.
//
// Note that ExcelWriter currently only handles
// string and int types —
// with strings being output as text labels, and
// ints as numeric values.
type ExcelWriter struct {
	*base
	prefix        string
	suffix        string
	builderByType map[reflect.Type]*excelBuilder
}

// NewExcelWriter returns a new ExcelWriter, using prefix
// and suffix when building its output filenames.
//
// See ExcelWriter (above) for output filename details.
func NewExcelWriter(prefix, suffix string) *ExcelWriter {
	w := ExcelWriter{
		base:          &base{},
		prefix:        prefix,
		suffix:        suffix,
		builderByType: make(map[reflect.Type]*excelBuilder),
	}
	return &w
}

func (w *ExcelWriter) register(x interface{}) (reflect.Type, error) {
	// Register with base writer.
	t, ok := w.base.register(x)
	if !ok {
		return t, nil
	}

	excel, err := newExcelBuilder(w.prefix + t.Name() + w.suffix + ".xlsx")
	if err != nil {
		return nil, err
	}
	w.builderByType[t] = excel

	h := convert(w.headersByType[t])
	err = excel.AddRow(h...)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func convert(list []string) []interface{} {
	var out []interface{}
	for _, x := range list {
		out = append(out, interface{}(x))
	}
	return out
}

// Write is called to persist records.
// Each record is written to an individual row
// in the corresponding output file, according to the
// type of the given record.
func (w *ExcelWriter) Write(x interface{}) error {
	t, err := w.register(x)
	if err != nil {
		return err
	}
	excel := w.builderByType[t]
	return excel.AddRow(excelValuesFrom(x)...)
}

// Close the writer, ensuring all files are saved.
func (w *ExcelWriter) Close() error {
	var rerr error
	for _, excel := range w.builderByType {
		err := excel.Save()
		if err != nil {
			rerr = err
		}
	}
	return rerr
}

// Cancel should be called in the event of an error occurring.
func (w *ExcelWriter) Cancel() error {
	// No clean up needed.
	return nil
}

// func excelHeaders(x interface{}) []interface{} {
// 	var out []interface{}
// 	reflectStructFields(x, func(name string, t reflect.Type, tag string) {
// 		tag = firstTagValue(tag)
// 		out = append(out, tag)
// 	})
// 	return out
// }

func excelValuesFrom(x interface{}) []interface{} {
	// TODO This is badly named, it's not just used by ExcelWriter.
	var out []interface{}
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		switch t.Kind() {
		case reflect.String, reflect.Int:
			out = append(out, v)
		default:
			m := fmt.Sprintf("Unknown type: %v", v)
			panic(m)
		}
		// out = append(out, v)
	})
	return out
}
