package peanut

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	// Import Sqlite db driver.
	_ "github.com/mattn/go-sqlite3"
)

// TODO Add an option to allow different insert modes (default/ignore/update).

type SQLiteWriter struct {
	*base
	tmpFilename  string                     // tmpFilename is the filename used by the temp file.
	dstFilename  string                     // dstFilename is the final destination filename.
	insertByType map[reflect.Type]*sql.Stmt // insertByType holds prepared INSERT statements.
	db           *sql.DB                    // db is the database instance.
}

// TODO Can we unify/simplify the constructors? Use pattern instead of prefix/suffix maybe? (not here, but for others)

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

func (w *SQLiteWriter) createDDL(t reflect.Type) string {

	// Create table using type name.
	ddl := "CREATE TABLE " + t.Name() + " (\n"

	// List of DDL statements to build the table definition.
	var ddlLines []string
	// List of primary keys.
	var pks []string

	hdrs := w.headersByType[t]
	typs := w.typesByType[t]
	tags := w.tagsByType[t]

	for i := 0; i < len(typs); i++ {
		// Column name.
		col := "\t" + hdrs[i] + " "

		// Column datatype.
		switch typs[i].Kind() {
		case reflect.String:
			col += "TEXT"
		case reflect.Int:
			col += "INTEGER"
		default:
			m := fmt.Sprintf("Unknown type: %v", t.Name())
			panic(m)
		}

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
	s := "INSERT OR IGNORE INTO " + t.Name() + " ("
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

func (w *SQLiteWriter) Close() error {
	var rerr error

	// TODO

	err := w.close()
	if err != nil {
		rerr = err
	}

	err = os.Rename(w.tmpFilename, w.dstFilename)
	if err != nil {
		log.Printf("Error %s", err)
		rerr = err
	}

	return rerr
}

func (w *SQLiteWriter) close() error {
	var rerr error

	// TODO

	for _, stmt := range w.insertByType {
		var cerr error
		var err error

		err = stmt.Close()
		if err != nil {
			log.Printf("Error %s", err)
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

func (w *SQLiteWriter) Cancel() error {
	var rerr error

	// TODO

	rerr = w.close()

	err := os.Remove(w.tmpFilename)
	if err != nil {
		log.Printf("Error %s", err)
		rerr = err
	}

	return rerr
}
