package peanut

// JSONLWriter
// - See https://en.wikipedia.org/wiki/JSON_streaming
// - To futz about with tags/etc, see https://stackoverflow.com/questions/42546519/how-do-i-dynamically-change-the-structs-json-tag

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// TODO This is based upon CSVWriter. Complete this...

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
//
// Note that JSONLWriter currently only handles
// string and int types:
// strings are output as strings, and ints as numbers.
type JSONLWriter struct {
	*base
	prefix      string
	suffix      string
	jsonlByType map[reflect.Type]*jsonlBuilder
}

// NewJSONLWriter returns a new JSONLWriter, using prefix
// and suffix when building its output filenames.
//
// See JSONLWriter (above) for output filename details.
func NewJSONLWriter(prefix, suffix string) *JSONLWriter {
	w := JSONLWriter{
		base:        &base{},
		prefix:      prefix,
		suffix:      suffix,
		jsonlByType: make(map[reflect.Type]*jsonlBuilder),
	}
	return &w
}

type jsonlBuilder struct {
	filename string
	file     *os.File
	bw       *bufio.Writer
	enc      *json.Encoder
}

func (w *JSONLWriter) register(x interface{}) error {
	// Register with base writer.
	ok := w.base.register(x)
	if !ok {
		return nil
	}

	i := len(w.types) - 1
	t := w.types[i]

	// log.Printf("Setting up jsonl.Writer for %s", t.Name())

	name := w.prefix + t.Name() + w.suffix + ".jsonl"
	// file, err := os.Create(name)
	file, err := ioutil.TempFile("", "atomic-")
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	bw := bufio.NewWriter(file)
	enc := json.NewEncoder(bw)
	w.jsonlByType[t] = &jsonlBuilder{filename: name, file: file, bw: bw, enc: enc}

	// err = cw.Write(w.typeHeaders[i])
	// if err != nil {
	// 	log.Printf("Error %s", err)
	// 	// w.Destroy()
	// 	return err
	// }
	return nil
}

// Write is called to persist records.
// Each record is written to an individual row
// in the corresponding output file, according to the
// type of the given record.
func (w *JSONLWriter) Write(x interface{}) error {
	err := w.register(x)
	if err != nil {
		return err
	}
	t := baseType(x)
	// log.Printf("WriteRecord for %s", t.Name())

	// TODO
	// panic("unimplemented")

	enc := w.jsonlByType[t].enc
	return enc.Encode(mapValues(x))
}

func mapValues(x interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	reflectStructValues(x, func(name string, t reflect.Type, v interface{}, tag string) {
		tag = firstTagValue(tag)
		switch t.Kind() {
		case reflect.String:
			out[tag] = v.(string)
			// out = append(out, v.(string))
		case reflect.Int:
			out[tag] = v.(int)
			// out = append(out, strconv.Itoa(v.(int)))
		default:
			m := fmt.Sprintf("Unknown type: %v", v)
			panic(m)
		}
	})
	return out
}

// Close flushes all buffers and writers,
// and closes the output files.
func (w *JSONLWriter) Close() error {

	// // TODO
	// panic("unimplemented")

	var rerr error
	for _, c := range w.jsonlByType {
		var cerr error
		var err error
		// c.csvw.Flush()
		// err = c.csvw.Error()
		// if err != nil {
		// 	log.Printf("Error %s", err)
		// 	cerr = err
		// }

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
func (w *JSONLWriter) Cancel() error {

	// TODO
	// panic("unimplemented")

	var rerr error
	for _, c := range w.jsonlByType {
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
