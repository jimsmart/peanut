package peanut

import (
	"database/sql"
	"os"
	"reflect"
	"strings"

	// Import Sqlite db driver.
	_ "github.com/mattn/go-sqlite3"
)

// TODO(js) Add an option to allow different insert modes (default/ignore/update).

// TODO(js) We should not be silently overwriting things. We should return an error somehow.

// SQLiteWriter writes records to an SQLite database,
// writing each record type to an individual table
// automatically.
//
// During writing, the database file is held in a
// temporary location, and only moved into its
// final destination during a successful Close operation.
//
// Note that if an existing database with the same filename
// already exists at the given output location,
// it will be silently overwritten.
//
// The caller must call Close on successful completion
// of all writing, to ensure proper cleanup, and the
// relocation of the database from its temporary
// location during writing, to its final output.
//
// In the event of an error or cancellation, the
// caller must call Cancel before quiting, to ensure
// closure and cleanup of any partially written data.
//
// SQLiteWriter supports additional tag values to denote the primary key:
//  type Shape struct {
//  	ShapeID  string `peanut:"shape_id,pk"`
//  	Name     string `peanut:"name"`
//  	NumSides int    `peanut:"num_sides"`
//  }
//
//  type Color struct {
//  	ColorID string `peanut:"color_id,pk"`
//  	Name    string `peanut:"name"`
//  	RBG     string `peanut:"rgb"`
//  }
// Compound primary keys are also supported.
//
// SQLiteWriter has no support for foreign keys, indexes, etc.
type SQLiteWriter struct {
	*base
	tmpFilename  string                     // tmpFilename is the filename used by the temp file.
	dstFilename  string                     // dstFilename is the final destination filename.
	insertByType map[reflect.Type]*sql.Stmt // insertByType holds prepared INSERT statements.
	db           *sql.DB                    // db is the database instance.
}

// TODO(js) Can we unify/simplify the constructors? Use pattern instead of prefix/suffix maybe? (not here, but for others)

// NewSQLiteWriter returns a new SQLiteWriter,
// using the given filename + ".sqlite" as its final output location.
func NewSQLiteWriter(filename string) *SQLiteWriter {
	w := SQLiteWriter{
		base:         &base{},
		dstFilename:  filename + ".sqlite",
		insertByType: make(map[reflect.Type]*sql.Stmt),
	}
	return &w
}

func (w *SQLiteWriter) register(x interface{}) (reflect.Type, error) {
	// Register with base writer.
	t, ok := w.base.register(x)
	if !ok {
		return t, nil
	}

	// Lazy init of database.
	if w.db == nil {

		filename, err := randomTempFilename("peanut-", ".sqlite")
		if err != nil {
			return nil, err
		}

		// log.Printf("Creating SQLite db %s", filename)
		db, err := sql.Open("sqlite3", filename)
		if err != nil {
			return nil, err
		}
		w.db = db
		w.tmpFilename = filename
	}

	// log.Printf("Setting up SQLite table for %s", t.Name())

	ddl := w.createDDL(t)
	// log.Println("DDL:", ddl)

	// Execute DDL to create table.
	_, err := w.db.Exec(ddl)
	if err != nil {
		return nil, err
	}

	insert := w.createInsert(t)
	// log.Println("Insert:", insert)

	// Create and cache prepared statement.
	stmt, err := w.db.Prepare(insert)
	if err != nil {
		return nil, err
	}
	w.insertByType[t] = stmt

	return t, nil
}

var kindToDBType = map[reflect.Kind]string{
	reflect.String:  "TEXT",
	reflect.Bool:    "BOOLEAN",
	reflect.Float64: "REAL",
	reflect.Float32: "REAL",
	reflect.Int8:    "INT8",
	reflect.Int16:   "INT16",
	reflect.Int32:   "INT32",
	reflect.Int64:   "INT64",
	reflect.Int:     "INT64",
	reflect.Uint8:   "UNSIGNED INT8",
	reflect.Uint16:  "UNSIGNED INT16",
	reflect.Uint32:  "UNSIGNED INT32",
	reflect.Uint64:  "UNSIGNED INT64",
	reflect.Uint:    "UNSIGNED INT64",
}

func (w *SQLiteWriter) createDDL(t reflect.Type) string {

	// Create table using type name.
	ddl := "CREATE TABLE \"" + t.Name() + "\" (\n"

	// List of DDL statements to build the table definition.
	var ddlLines []string
	// List of primary keys.
	var pks []string

	hdrs := w.headersByType[t]
	typs := w.typesByType[t]
	tags := w.tagsByType[t]

	for i := 0; i < len(typs); i++ {
		// TODO(js) We should quote this appropriately, to handle reserved words.
		// Column name.
		col := "\t\"" + hdrs[i] + "\" "

		// Column datatype.
		col += kindToDBType[typs[i].Kind()]

		// Column constraints.
		col += " NOT NULL"

		// Add DDL line to list.
		ddlLines = append(ddlLines, col)

		// Handle primary key tag.
		if secondTagValue(tags[i]) == "pk" {
			// Add column name to primary key list.
			pks = append(pks, hdrs[i])
		}
	}

	// Primary key.
	if len(pks) > 0 {
		pk := "PRIMARY KEY ("
		pk += strings.Join(pks, ", ")
		pk += ")"
		ddlLines = append(ddlLines, pk)
	}

	// Join the lines of DDL together.
	ddl += strings.Join(ddlLines, ",\n")

	ddl += "\n)"
	return ddl
}

func (w *SQLiteWriter) createInsert(t reflect.Type) string {
	// TODO Add an option to allow different insert modes (default/ignore/update).
	s := "INSERT OR IGNORE INTO \"" + t.Name() + "\" ("
	hdrs := w.headersByType[t]
	s += strings.Join(hdrs, ",")
	s += ") VALUES ("
	q := make([]string, len(hdrs))
	for i := 0; i < len(q); i++ {
		q[i] = "?"
	}
	s += strings.Join(q, ",")
	s += ")"
	return s
}

// Write is called to persist records.
// Each record is written to an individual row
// in the corresponding table within the output database,
// according to the type of the given record.
func (w *SQLiteWriter) Write(x interface{}) error {
	t, err := w.register(x)
	if err != nil {
		return err
	}

	// log.Printf("WriteRecord for %s", t.Name())
	stmt := w.insertByType[t]
	_, err = stmt.Exec(excelValuesFrom(x)...)
	return err
}

// Close cleans up all used resources,
// closes the database connection,
// and moves the database to its final location.
func (w *SQLiteWriter) Close() error {
	var rerr error

	// TODO(js) We should make lists of errors.

	err := w.close()
	if err != nil {
		rerr = err
	}

	err = os.Rename(w.tmpFilename, w.dstFilename)
	if err != nil {
		rerr = err
	}

	return rerr
}

func (w *SQLiteWriter) close() error {
	var rerr error

	// TODO(js) We should make lists of errors.

	for _, stmt := range w.insertByType {
		var cerr error
		var err error

		err = stmt.Close()
		if err != nil {
			cerr = err
		}

		if cerr != nil {
			rerr = cerr
		}
	}

	err := w.db.Close()
	if err != nil {
		rerr = err
	}

	return rerr
}

// Cancel should be called in the event of an error occurring,
// to properly close any used resources,
// and delete the partially written database from its temporary location.
func (w *SQLiteWriter) Cancel() error {
	var rerr error

	// TODO(js) We should make lists of errors.

	rerr = w.close()

	err := os.Remove(w.tmpFilename)
	if err != nil {
		rerr = err
	}

	return rerr
}
