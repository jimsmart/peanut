// Package peanut writes data structs in a variety of formats.
// Its primary job is to provide a single consistent interface
// for easy, ceremony-free persistence of struct data.
//
// Currently supported formats are CSV and Excel files. Additional
// writers are also provided to assist with testing and debugging.
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
// fields. e.g.
//  type ParentRecord struct {
// 	ParentID string `peanut:"parent_id"`
// 	Name     string `peanut:"name"`
// 	Counter  int    `peanut:"counter"`
//  }
//
//  type ChildRecord struct {
// 	ChildID  string `peanut:"child_id"`
// 	Name     string `peanut:"name"`
// 	ParentID string `peanut:"parent_id"`
//  }
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
// 	// Handle error...
//
// 	y := &ChildRecord{
// 		ChildID:  "cid1",
// 		Name:     "bar",
// 		ParentID: "pid1",
// 	}
// 	err = w.Write(y)
// 	// Handle error...
//
// 	z := &ParentRecord{
// 		ParentID: "pid2",
// 		Name:     "baz",
// 		Counter:  456,
// 	}
// 	err = w.Write(z)
// 	// Handle error...
//
// When successfully completed:
//  err = w.Close()
//
// Or, to abort the whole operation in the event of an error or cancellation:
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
// but in different packages, e.g. package1.Foo and package2.Foo.
//
// Currently only types string and int are fully supported by all writers.
package peanut
