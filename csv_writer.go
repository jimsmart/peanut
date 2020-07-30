package peanut

import (
	"bufio"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// TODO(js) csv.Writer has an internal bufio.Writer, remove our usage of same.

var _ Writer = &CSVWriter{}

// CSVWriter writes records to CSV files, writing
// each record type to an individual CSV file automatically.
//
// Filenames for each corresponding record type are derived
// accordingly:
//  prefix + type.Name() + suffix + ".csv"
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
//
// Note that CSVWriter currently only handles
// string and int types,
// both of which are output as strings.
type CSVWriter struct {
	*base
	prefix        string
	suffix        string
	builderByType map[reflect.Type]*csvBuilder
}

// NewCSVWriter returns a new CSVWriter, using prefix
// and suffix when building its output filenames.
//
// See CSVWriter (above) for output filename details.
func NewCSVWriter(prefix, suffix string) *CSVWriter {
	w := CSVWriter{
		base:          &base{},
		prefix:        prefix,
		suffix:        suffix,
		builderByType: make(map[reflect.Type]*csvBuilder),
	}
	return &w
}

type csvBuilder struct {
	filename string
	file     *os.File
	bw       *bufio.Writer // TODO Remove this.
	csvw     *csv.Writer
}

func (w *CSVWriter) register(x interface{}) error {
	// Register with base writer.
	ok := w.base.register(x)
	if !ok {
		return nil
	}

	i := len(w.types) - 1
	t := w.types[i]

	// log.Printf("Setting up csv.Writer for %s", t.Name())

	name := w.prefix + t.Name() + w.suffix + ".csv"
	// file, err := os.Create(name)
	file, err := ioutil.TempFile("", "atomic-")
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	bw := bufio.NewWriter(file) // TODO Remove this.
	cw := csv.NewWriter(bw)
	w.builderByType[t] = &csvBuilder{filename: name, file: file, bw: bw, csvw: cw}

	err = cw.Write(w.typeHeaders[i])
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	return nil
}

// Write is called to persist records.
// Each record is written to an individual row
// in the corresponding output file, according to the
// type of the given record.
func (w *CSVWriter) Write(x interface{}) error {
	err := w.register(x)
	if err != nil {
		return err
	}
	t := baseType(x)
	// log.Printf("WriteRecord for %s", t.Name())
	cw := w.builderByType[t].csvw
	return cw.Write(stringValues(x))
}

// Close flushes all buffers and writers,
// and closes the output files.
func (w *CSVWriter) Close() error {
	var rerr error
	for _, c := range w.builderByType {
		var cerr error
		var err error
		c.csvw.Flush()
		err = c.csvw.Error()
		if err != nil {
			log.Printf("Error %s", err)
			cerr = err
		}

		err = c.bw.Flush() // TODO Remove this.
		if err != nil {
			log.Printf("Error %s", err)
			cerr = err
		}

		// Chmod the file world-readable (ioutil.TempFile creates files with
		// mode 0600) before renaming.
		err = c.file.Chmod(0644)
		if err != nil {
			log.Printf("Error %s", err)
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
			log.Printf("Error %s", err)
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
			log.Printf("Error %s", err)
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
	var rerr error
	for _, c := range w.builderByType {
		var err error

		err = c.file.Close()
		if err != nil {
			log.Printf("Error %s", err)
			rerr = err
		}

		err = os.Remove(c.file.Name())
		if err != nil {
			log.Printf("Error %s", err)
			rerr = err
		}
	}
	return rerr
}
