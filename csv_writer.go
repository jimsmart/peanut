package peanut

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"reflect"
)

var _ Writer = &CSVWriter{}

// CSVWriter writes records to CSV files, writing
// each record type to an individual CSV file automatically.
//
// Filenames for each corresponding record type are derived
// accordingly:
//  prefix + type.Name() + suffix + extension
//
// Where extension is ".csv" or ".tsv" accordingly.
//
// The first row of resulting CSV file(s) will contain
// headers using names extracted from the struct's
// field tags. Records' fields are written in the order
// that they appear within the struct.
//
// The caller must call Close on successful completion
// of all writing, to ensure buffers are flushed and
// files are properly written to disk.
//
// In the event of an error or cancellation, the
// caller must call Cancel before quiting, to ensure
// closure and cleanup of any partially written files.
type CSVWriter struct {
	*base
	prefix        string
	suffix        string
	extension     string
	comma         rune
	builderByType map[reflect.Type]*csvBuilder
}

// NewCSVWriter returns a new CSVWriter, using prefix
// and suffix when building its output filenames,
// and using ".csv" file extension with comma ',' as a field separator.
//
// See CSVWriter (above) for output filename details.
func NewCSVWriter(prefix, suffix string) *CSVWriter {
	w := CSVWriter{
		base:          &base{},
		prefix:        prefix,
		suffix:        suffix,
		extension:     ".csv",
		comma:         ',',
		builderByType: make(map[reflect.Type]*csvBuilder),
	}
	return &w
}

// NewTSVWriter returns a new CSVWriter configured to write
// TSV files, using prefix and suffix when building its output filenames,
// and using ".tsv" file extension with tab '\t' as a field separator.
//
// See CSVWriter (above) for output filename details.
func NewTSVWriter(prefix, suffix string) *CSVWriter {
	w := NewCSVWriter(prefix, suffix)
	w.extension = ".tsv"
	w.comma = '\t'
	return w
}

type csvBuilder struct {
	filename string
	file     *os.File
	csvw     *csv.Writer
}

func (w *CSVWriter) register(x interface{}) (reflect.Type, error) {
	// Register with base writer.
	t, ok := w.base.register(x)
	if !ok {
		return t, nil
	}
	if err := allFieldsSupportedKinds(x); err != nil {
		return nil, err
	}
	if len(w.base.tagsByType[t]) == 0 {
		return t, nil
	}

	// log.Printf("Setting up csv.Writer for %s", t.Name())

	name := w.prefix + t.Name() + w.suffix + w.extension
	file, err := ioutil.TempFile("", "atomic-")
	if err != nil {
		return nil, err
	}
	cw := csv.NewWriter(file)
	cw.Comma = w.comma
	w.builderByType[t] = &csvBuilder{filename: name, file: file, csvw: cw}

	err = cw.Write(w.headersByType[t])
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Write is called to persist records.
// Each record is written to an individual row
// in the corresponding output file, according to the
// type of the given record.
func (w *CSVWriter) Write(x interface{}) error {
	if w.closed {
		return ErrClosedWriter
	}
	t, err := w.register(x)
	if err != nil {
		return err
	}
	if len(w.base.tagsByType[t]) == 0 {
		return nil
	}
	// log.Printf("WriteRecord for %s", t.Name())
	cw := w.builderByType[t].csvw
	return cw.Write(stringValues(x))
}

// Close flushes all buffers and writers,
// and closes the output files.
//
// Calling Close after a previous call to
// Cancel is safe, and always results in a no-op.
func (w *CSVWriter) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true
	var rerr error
	for _, c := range w.builderByType {
		var cerr error
		var err error
		c.csvw.Flush()
		err = c.csvw.Error()
		if err != nil {
			cerr = err
		}

		// Chmod the file world-readable (ioutil.TempFile creates files with
		// mode 0600) before renaming.
		err = c.file.Chmod(0644)
		if err != nil {
			cerr = err
		}

		// fsync(2) after fchmod(2) orders writes as per
		// https://lwn.net/Articles/270891/. Can be skipped for performance
		// for idempotent applications (which only ever atomically write new
		// files and tolerate file loss) on an ordered file systems. ext3,
		// ext4, XFS, Btrfs, ZFS are ordered by default.
		c.file.Sync()

		err = c.file.Close()
		if err != nil {
			cerr = err
		}

		if cerr != nil {
			rerr = cerr
			// // Best effort cleanup.
			// os.Remove(c.file.Name())
			continue
		}

		err = os.Rename(c.file.Name(), c.filename)
		if err != nil {
			cerr = err
		}

		if cerr != nil {
			rerr = cerr
		}
	}
	return rerr
}

// Cancel should be called in the event of an error occurring,
// to properly close and delete any partially written files.
func (w *CSVWriter) Cancel() error {
	if w.closed {
		return nil
	}
	w.closed = true
	var rerr error
	for _, c := range w.builderByType {
		var err error

		err = c.file.Close()
		if err != nil {
			rerr = err
		}

		err = os.Remove(c.file.Name())
		if err != nil {
			rerr = err
		}
	}
	return rerr
}
