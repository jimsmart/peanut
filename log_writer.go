package peanut

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

// TODO(js) This could have a multi-line mode, where each field is logged to its own line.

var _ Writer = &LogWriter{}

// LogWriter writes records to a log.Logger.
//
// If Logger is nil at runtime, a new log.Logger
// will be created when needed, writing to os.Stderr.
type LogWriter struct {
	*base
	Logger  *log.Logger
	Verbose bool
}

func (w *LogWriter) register(x interface{}) (reflect.Type, error) {
	if w.base == nil {
		w.base = &base{}
		if w.Logger == nil {
			w.Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
		}
	}
	t, ok := w.base.register(x)
	if !ok {
		return t, nil
	}
	if err := allFieldsSupportedKinds(x); err != nil {
		return nil, err
	}
	return t, nil
}

// Write is called to persist records.
func (w *LogWriter) Write(x interface{}) error {
	t, err := w.register(x)
	if err != nil {
		return err
	}
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

// Close should be called after successfully writing records.
func (w *LogWriter) Close() error {
	if w.Verbose {
		w.Logger.Println("Called LogWriter.Close")
	}
	return nil
}

// Cancel should be called in the event of an error occurring.
func (w *LogWriter) Cancel() error {
	if w.Verbose {
		w.Logger.Println("Called LogWriter.Cancel")
	}
	return nil
}
