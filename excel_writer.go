package peanut

import (
	"log"
	"reflect"
)

var _ Writer = &ExcelWriter{}

// ExcelWriter is a type of peanut.Writer that writes Excel files.
type ExcelWriter struct {
	*writer
	prefix      string
	suffix      string
	excelByType map[reflect.Type]*excelBuilder
}

// NewExcelWriter returns a new ExcelWriter.
func NewExcelWriter(prefix, suffix string) *ExcelWriter {
	w := ExcelWriter{
		writer:      &writer{},
		prefix:      prefix,
		suffix:      suffix,
		excelByType: make(map[reflect.Type]*excelBuilder),
	}
	return &w
}

func (w *ExcelWriter) initialise(x interface{}) error {
	// Register with base writer.
	ok := w.init(x)
	if !ok {
		return nil
	}

	i := len(w.types) - 1
	t := w.types[i]

	excel, err := newExcelBuilder(w.prefix + t.Name() + w.suffix + ".xlsx")
	if err != nil {
		return err
	}
	w.excelByType[t] = excel

	h := convert(w.headers[i])
	err = excel.AddRow(h...)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	return nil
}

func convert(list []string) []interface{} {
	var out []interface{}
	for _, x := range list {
		out = append(out, interface{}(x))
	}
	return out
}

// Write is called to persist records.
func (w *ExcelWriter) Write(x interface{}) error {
	err := w.initialise(x)
	if err != nil {
		return err
	}
	t := baseType(x)
	excel := w.excelByType[t]
	return excel.AddRow(excelValuesFrom(x)...)
}

// Close the writer, ensuring all files are saved.
func (w *ExcelWriter) Close() error {
	var rerr error
	for _, excel := range w.excelByType {
		err := excel.Save()
		if err != nil {
			log.Printf("Error %s", err)
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
	var out []interface{}
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		out = append(out, v)
	})
	return out
}
