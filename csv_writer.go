package peanut

import (
	"bufio"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

var _ Writer = &CSVWriter{}

// CSVWriter writes records to CSV files.
type CSVWriter struct {
	*writer
	prefix    string
	suffix    string
	csvByType map[reflect.Type]*csvBuilder
}

// NewCSVWriter returns a new CSVWriter.
func NewCSVWriter(prefix, suffix string) *CSVWriter {
	w := CSVWriter{
		writer:    &writer{},
		prefix:    prefix,
		suffix:    suffix,
		csvByType: make(map[reflect.Type]*csvBuilder),
	}
	return &w
}

type csvBuilder struct {
	filename string
	file     *os.File
	bw       *bufio.Writer
	csvw     *csv.Writer
}

func (w *CSVWriter) initialise(x interface{}) error {
	// Register with base writer.
	ok := w.init(x)
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
	bw := bufio.NewWriter(file)
	cw := csv.NewWriter(bw)
	w.csvByType[t] = &csvBuilder{filename: name, file: file, bw: bw, csvw: cw}

	err = cw.Write(w.headers[i])
	if err != nil {
		log.Printf("Error %s", err)
		// w.Destroy()
		return err
	}
	return nil
}

// Write is called to persist records.
func (w *CSVWriter) Write(x interface{}) error {
	err := w.initialise(x)
	if err != nil {
		return err
	}
	t := baseType(x)
	// log.Printf("WriteRecord for %s", t.Name())
	cw := w.csvByType[t].csvw
	return cw.Write(stringValues(x))
}

// Close flushes all buffers and writers,
// and closes the output files.
func (w *CSVWriter) Close() error {
	var rerr error
	for _, c := range w.csvByType {
		var cerr error
		var err error
		c.csvw.Flush()
		err = c.csvw.Error()
		if err != nil {
			log.Printf("Error %s", err)
			cerr = err
		}

		err = c.bw.Flush()
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
	for _, c := range w.csvByType {
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
