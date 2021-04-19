package peanut_test

import (
	"database/sql"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
	"github.com/jimsmart/schema"

	// Import Sqlite db driver.
	_ "github.com/mattn/go-sqlite3"
)

type tableResults struct {
	columns []string
	types   []string
	pks     []string
	data    [][]string
}

var _ = Describe("SQLiteWriter", func() {

	newFn := func(suffix string) peanut.Writer {
		w := peanut.NewSQLiteWriter("./test/output" + suffix)
		return w
	}

	expectedOutput := map[string]*tableResults{
		"Foo": {
			columns: []string{"foo_string", "foo_int"},
			types:   []string{"TEXT", "INT64"},
			pks:     []string{"foo_string"},
			data: [][]string{
				{"test 1", "1"},
				{"test 2", "2"},
				{"test 3", "3"},
			},
		},
		"Bar": {
			columns: []string{"bar_int", "bar_string"},
			types:   []string{"INT64", "TEXT"},
			pks:     []string{"bar_int", "bar_string"},
			data: [][]string{
				{"1", "test 1"},
				{"2", "test 2"},
				{"3", "test 3"},
			},
		},
		"Baz": {
			columns: []string{"baz_string", "baz_bool", "baz_float32", "baz_float64", "baz_int", "baz_int8", "baz_int16", "baz_int32", "baz_int64", "baz_uint", "baz_uint8", "baz_uint16", "baz_uint32", "baz_uint64"},
			types:   []string{"TEXT", "BOOLEAN", "REAL", "REAL", "INT64", "INT8", "INT16", "INT32", "INT64", "UNSIGNED INT64", "UNSIGNED INT8", "UNSIGNED INT16", "UNSIGNED INT32", "UNSIGNED INT64"},
			pks:     []string{"baz_string"},
			data: [][]string{
				{"test 1", "true", "1.234", "9.876", "-12345", "-8", "-16", "-32", "-64", "12345", "8", "16", "32", "64"},
			},
		},
	}

	AfterEach(func() {
		os.Remove("./test/output-sequential.sqlite")
		os.Remove("./test/output-interleave.sqlite")
	})

	It("should write the correct data when sequential structs are written", func() {
		w := newFn("-sequential")

		testWritesAndCloseSequential(w)

		output1, err := readSQLite("./test/output-sequential.sqlite")
		Expect(err).To(BeNil())
		Expect(output1).To(Equal(expectedOutput))
	})

	It("should write the correct data when interleaved structs are written", func() {
		w := newFn("-interleave")

		testWritesAndCloseSequential(w)

		output1, err := readSQLite("./test/output-interleave.sqlite")
		Expect(err).To(BeNil())
		Expect(output1).To(Equal(expectedOutput))
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-cancel.sqlite").ToNot(BeAnExistingFile())
	})

	It("should return an error when Write is called after Close", func() {
		w := newFn("-close-write")

		testWriteAfterClose(w)

		Expect("./test/output-Foo-close-write.sqlite").ToNot(BeAnExistingFile())
	})

	It("should return an error when the path is bad", func() {
		w := peanut.NewSQLiteWriter("./no-such-location/output-bogus")

		err := w.Write(testOutputFoo[0])
		Expect(err).To(BeNil())

		err = w.Close()
		Expect(err).ToNot(BeNil())
	})

	Context("when given a struct with an unsupported field type", func() {

		It("should return an error with an informative message", func() {
			w := peanut.NewSQLiteWriter("./no-such-location/output-bogus")
			defer func() {
				err1 := w.Cancel()
				err2 := w.Close()
				Expect(err1).To(BeNil())
				Expect(err2).To(BeNil())
			}()

			err := w.Write(BadUnsupported{})
			Expect(err).ToNot(BeNil())

			// Expect error message to be informative.
			Expect(err.Error()).To(SatisfyAll(
				MatchRegexp(`slice`),          // type
				MatchRegexp("BytesField"),     // field name
				MatchRegexp("BadUnsupported"), // struct name
			))
		})
	})

})

func readSQLite(filename string) (map[string]*tableResults, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	out := make(map[string]*tableResults)

	tables, err := schema.Tables(db)
	if err != nil {
		return nil, err
	}
	for name, ct := range tables {
		var headers []string
		var dbtypes []string
		for i := range ct {
			headers = append(headers, ct[i].Name())
			dbtypes = append(dbtypes, ct[i].DatabaseTypeName())
		}

		pks, err := schema.PrimaryKey(db, "", name[1])
		if err != nil {
			return nil, err
		}

		data, err := readData(db, name[1])
		if err != nil {
			return nil, err
		}

		result := &tableResults{
			columns: headers,
			types:   dbtypes,
			pks:     pks,
			data:    data,
		}

		out[name[1]] = result
	}
	return out, nil
}

func readData(db *sql.DB, table string) ([][]string, error) {
	q := "SELECT * FROM " + table
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rawResults := make([][]byte, len(cols))
	results := make([]string, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range rawResults {
		// Put pointers to each string in the interface slice
		dest[i] = &rawResults[i]
	}

	var out [][]string
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		for i, raw := range rawResults {
			if raw != nil {
				results[i] = string(raw)
			}
		}
		row := make([]string, len(cols))
		copy(row, results)
		out = append(out, row)
	}

	return out, nil
}
