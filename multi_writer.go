package peanut

var _ Writer = &multiWriter{}

// multiWriter implements peanut.Writer
type multiWriter struct {
	writers []Writer
}

// MultiWriter creates a writer that duplicates its method calls to all the provided writers.
func MultiWriter(writers ...Writer) Writer {
	return &multiWriter{writers: writers}
}

func (mw *multiWriter) Write(x interface{}) error {
	for _, w := range mw.writers {
		err := w.Write(x)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mw *multiWriter) Close() error {
	var xerr error
	for _, w := range mw.writers {
		err := w.Close()
		if err != nil {
			xerr = err
		}
	}
	return xerr
}

func (mw *multiWriter) Cancel() error {
	var xerr error
	for _, w := range mw.writers {
		err := w.Cancel()
		if err != nil {
			xerr = err
		}
	}
	return xerr
}
