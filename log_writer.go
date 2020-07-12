package peanut

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

// TODO This could have a multi-line mode, where each field is logged to its own line.

var _ Writer = &LogWriter{}

// implements peanut.Writer
type LogWriter struct {
	*writer
	Logger  *log.Logger
	Verbose bool
}

func (w *LogWriter) initialise(x interface{}) error {
	if w.writer == nil {
		w.writer = &writer{}
		if w.Logger == nil {
			w.Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
		}
	}
	w.init(x)
	return nil
}

// Write is called to persist records.
func (w *LogWriter) Write(x interface{}) error {
	err := w.initialise(x)
	if err != nil {
		return err
	}
	t := baseType(x)
	// Build log message.
	n := t.Name()
	m := fmt.Sprintf("<%s>", n)
	// Concatenate field names and values.
	h := w.headersByType[t]
	v := stringValues(x)
	for i := range v {
		m += fmt.Sprintf(" %s:", h[i])
		if len(v[i]) > 0 {
			m += fmt.Sprintf(" %s", v[i])
		}
	}
	// Log it.
	w.Logger.Println(m)
	return nil
}

func (w *LogWriter) Close() error {
	if w.Verbose {
		w.Logger.Println("Called LogWriter.Close")
	}
	return nil
}

func (w *LogWriter) Cancel() error {
	if w.Verbose {
		w.Logger.Println("Called LogWriter.Cancel")
	}
	return nil
}

func stringValues(x interface{}) []string {
	var out []string
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		switch t.Kind() {
		case reflect.String:
			out = append(out, v.(string))
		case reflect.Int:
			out = append(out, strconv.Itoa(v.(int)))
		default:
			m := fmt.Sprintf("Unknown type: %v", v)
			panic(m)
		}
	})
	return out
}
