package peanut

// JSONLWriter
// - See http://jsonlines.org
// - See https://en.wikipedia.org/wiki/JSON_streaming
// - To futz about with tags/etc, see https://stackoverflow.com/questions/42546519/how-do-i-dynamically-change-the-structs-json-tag

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
)

var _ Writer = &JSONLWriter{}

// JSONLWriter writes records to JSON Lines files, writing
// each record type to an individual JSON Lines file automatically.
//
// Filenames for each corresponding record type are derived
// accordingly:
//  prefix + type.Name() + suffix + ".jsonl"
//
// The caller must call Close on successful completion
// of all writing, to ensure buffers are flushed and
// files are properly written to disk.
//
// In the event of an error or cancellation, the
// caller must call Cancel before quiting, to ensure
// closure and cleanup of any partially written files.
type JSONLWriter struct {
	*base
	prefix        string
	suffix        string
	builderByType map[reflect.Type]*jsonlBuilder
}

// NewJSONLWriter returns a new JSONLWriter, using prefix
// and suffix when building its output filenames.
//
// See JSONLWriter (above) for output filename details.
func NewJSONLWriter(prefix, suffix string) *JSONLWriter {
	w := JSONLWriter{
		base:          &base{},
		prefix:        prefix,
		suffix:        suffix,
		builderByType: make(map[reflect.Type]*jsonlBuilder),
	}
	return &w
}

type jsonlBuilder struct {
	filename string
	file     *os.File
	bw       *bufio.Writer
	enc      *json.Encoder
}

func (w *JSONLWriter) register(x interface{}) (reflect.Type, error) {
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

	// log.Printf("Setting up jsonl.Writer for %s", t.Name())

	name := w.prefix + t.Name() + w.suffix + ".jsonl"
	// file, err := os.Create(name)
	file, err := ioutil.TempFile("", "atomic-")
	if err != nil {
		return nil, err
	}
	bw := bufio.NewWriter(file)
	enc := json.NewEncoder(bw)
	enc.SetEscapeHTML(false)
	w.builderByType[t] = &jsonlBuilder{filename: name, file: file, bw: bw, enc: enc}
	return t, nil
}

// Write is called to persist records.
// Each record is written to an individual row
// in the corresponding output file, according to the
// type of the given record.
func (w *JSONLWriter) Write(x interface{}) error {
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
	enc := w.builderByType[t].enc
	return enc.Encode(mapValues(x))
}

// Close flushes all buffers and writers,
// and closes the output files.
//
// Calling Close after a previous call to
// Cancel is safe, and always results in a no-op.
func (w *JSONLWriter) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true

	var rerr error
	for _, c := range w.builderByType {
		var cerr error
		var err error
		err = c.bw.Flush()
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
func (w *JSONLWriter) Cancel() error {
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
