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

type SQLiteWriter struct {
	*base
	tmpFilename string
	dstFilename string
	// suffix        string
	// builderByType map[reflect.Type]*csvBuilder
	insertByType map[reflect.Type]*sql.Stmt
	db           *sql.DB
}

// TODO Can we unify/simplify the constructors? Use pattern instead of prefix/suffix maybe?

// TODO Prefix should likely just be filename here?

func NewSQLiteWriter(prefix string) *SQLiteWriter {
	w := SQLiteWriter{
		base:        &base{},
		dstFilename: prefix + ".sqlite",
		// suffix:        suffix,
		// builderByType: make(map[reflect.Type]*csvBuilder),
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

		log.Printf("Creating SQLite db %s", filename)
		db, err := sql.Open("sqlite3", filename)
		if err != nil {
			return nil, err
		}
		w.db = db
		w.tmpFilename = filename
	}

	log.Printf("Setting up SQLite table for %s", t.Name())

	ddl := w.createDDL(t)

	log.Println("DDL:", ddl)

	// Execute DDL to create table.
	_, err := w.db.Exec(ddl)
	if err != nil {
		return nil, err
	}

	insert := w.createInsert(t)

	log.Println("Insert:", insert)
	stmt, err := w.db.Prepare(insert)
	if err != nil {
		return nil, err
	}
	w.insertByType[t] = stmt

	return t, nil
}

func (w *SQLiteWriter) createDDL(t reflect.Type) string {
	ddl := "CREATE TABLE " + t.Name() + " (\n"

	var ddlLines []string
	var pks []string

	hdrs := w.headersByType[t]
	typs := w.typesByType[t]
	tags := w.tagsByType[t]
	for i := 0; i < len(typs); i++ {
		col := "\t" + hdrs[i] + " "
		switch typs[i].Kind() {
		case reflect.String:
			col += "TEXT"
		case reflect.Int:
			col += "INTEGER"
		default:
			m := fmt.Sprintf("Unknown type: %v", t.Name())
			panic(m)
		}

		col += " NOT NULL"

		ddlLines = append(ddlLines, col)

		if secondTagValue(tags[i]) == "pk" {
			pks = append(pks, hdrs[i])
		}
	}

	if len(pks) > 0 {
		pk := "PRIMARY KEY ("
		pk += strings.Join(pks, ", ")
		pk += ")"
		ddlLines = append(ddlLines, pk)
	}

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
	first := true
	for i := 0; i < len(hdrs); i++ {
		if !first {
			s += ","
		}
		first = false
		s += "?"
	}
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

	//

	// for _, c := range w.builderByType {
	// 	var cerr error
	// 	var err error
	// 	c.csvw.Flush()
	// 	err = c.csvw.Error()
	// 	if err != nil {
	// 		log.Printf("Error %s", err)
	// 		cerr = err
	// 	}

	// 	// Chmod the file world-readable (ioutil.TempFile creates files with
	// 	// mode 0600) before renaming.
	// 	err = c.file.Chmod(0644)
	// 	if err != nil {
	// 		log.Printf("Error %s", err)
	// 		cerr = err
	// 	}

	// 	// fsync(2) after fchmod(2) orders writes as per
	// 	// https://lwn.net/Articles/270891/. Can be skipped for performance
	// 	// for idempotent applications (which only ever atomically write new
	// 	// files and tolerate file loss) on an ordered file systems. ext3,
	// 	// ext4, XFS, Btrfs, ZFS are ordered by default.
	// 	c.file.Sync()

	// 	err = c.file.Close()
	// 	if err != nil {
	// 		log.Printf("Error %s", err)
	// 		cerr = err
	// 	}

	// 	if cerr != nil {
	// 		rerr = cerr
	// 		// // Best effort cleanup.
	// 		// os.Remove(c.file.Name())
	// 		continue
	// 	}

	// 	err = os.Rename(c.file.Name(), c.filename)
	// 	if err != nil {
	// 		log.Printf("Error %s", err)
	// 		cerr = err
	// 	}

	// 	if cerr != nil {
	// 		rerr = cerr
	// 	}
	// }
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

	//

	// for _, c := range w.builderByType {
	// 	var err error

	// 	err = c.file.Close()
	// 	if err != nil {
	// 		log.Printf("Error %s", err)
	// 		rerr = err
	// 	}

	// 	err = os.Remove(c.file.Name())
	// 	if err != nil {
	// 		log.Printf("Error %s", err)
	// 		rerr = err
	// 	}
	// }
	return rerr
}
