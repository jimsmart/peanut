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
			columns: []string{"foo_string1", "foo_int1"},
			types:   []string{"TEXT", "INTEGER"},
			pks:     []string{"foo_string1"},
			data: [][]string{
				{"test 1", "1"},
				{"test 2", "2"},
				{"test 3", "3"},
			},
		},
		"Bar": {
			columns: []string{"bar_int2", "bar_string2"},
			types:   []string{"INTEGER", "TEXT"},
			pks:     []string{"bar_int2", "bar_string2"},
			data: [][]string{
				{"1", "test 1"},
				{"2", "test 2"},
				{"3", "test 3"},
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
