// Package peanut writes data structs in a variety of formats.
// Its primary purpose is to provide a single consistent interface
// for easy, ceremony-free persistence of struct data.
//
// Each distinct struct type sent to Write() is written to an individual file/table,
// each named according to the name of the struct.
// Field/column names in each file/table are derived from struct tags.
//
// Currently supported formats are CSV, TSV, Excel (.xlsx), JSON Lines (JSONL), and SQLite.
// Additional writers are also provided to assist with testing and debugging.
// Mutiple writers can be combined using MultiWriter.
//
// All writers have the same basic interface: a Write method,
// that can take any appropriately tagged struct; a Close method,
// which should be called to successfully complete writing; and
// a Cancel method, which should be called to clean-up in the event
// of an error or cancellation during writing.
//
// Struct Tagging
//
// Structs to be used with peanut must have appropriately tagged
// fields, for example:
//  type ParentRecord struct {
//  	ParentID string `peanut:"parent_id"`
//  	Name     string `peanut:"name"`
//  	Counter  int    `peanut:"counter"`
//  }
//
//  type ChildRecord struct {
//  	ChildID  string `peanut:"child_id"`
//  	Name     string `peanut:"name"`
//  	ParentID string `peanut:"parent_id"`
//  }
//
// All writers use the same tags. Fields without tags do not get written as output.
//
// Usage
//
// First create a writer, for example:
//  w := peanut.NewCSVWriter("/some/path/my-data-", "")
//
// Next, write some records to it:
// 	x := &ParentRecord{
// 		ParentID: "pid1",
// 		Name:     "foo",
// 		Counter:  123,
// 	}
// 	err := w.Write(x)
// 	// ...
//
// 	y := &ChildRecord{
// 		ChildID:  "cid1",
// 		Name:     "bar",
// 		ParentID: "pid1",
// 	}
// 	err = w.Write(y)
// 	// ...
//
// 	z := &ParentRecord{
// 		ParentID: "pid2",
// 		Name:     "baz",
// 		Counter:  456,
// 	}
// 	err = w.Write(z)
// 	// ...
//
// When successfully completed:
//  err = w.Close()
//
//  // Output files will be:
//  // /some/path/my-data-ParentRecord.csv
//  // /some/path/my-data-ChildRecord.csv
//
// Or, to abort the whole operation in the event of an error or cancellation while writing records:
//  err = w.Cancel()
//
// MultiWriter
//
// Multiple writers can be combined using MultiWriter:
//  w1 := peanut.NewCSVWriter("/some/path/my-data-", "")
//  w2 := peanut.NewExcelWriter("/some/path/my-data-", "")
//  w3 := &peanut.LogWriter{}
//  w := peanut.MultiWriter(w1, w2, w3)
// Here w will write records to CSV, Excel and a logger.
//
// Limitations
//
// Behaviour is undefined for types with the same name
// but in different packages, such as package1.Foo and package2.Foo.
//
// Supported datatypes for struct fields: string, bool, float32, float64,
// int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64.
//
// Pointer following and nested structs are currently unsupported.
package peanut
