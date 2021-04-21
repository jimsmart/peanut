// Package peanut writes tagged data structs to disk in a variety of formats.
// Its primary purpose is to provide a single consistent interface
// for easy, ceremony-free persistence of record-based struct data.
//
// Each distinct struct type is written to an individual file (or table),
// automatically created, each named according to the name of the struct.
// Field/column names in each file/table are derived from struct tags.
// All writers use the same tags.
//
// Currently supported formats are CSV, TSV, Excel (.xlsx), JSON Lines (JSONL), and SQLite.
// Additional writers are also provided to assist with testing and debugging.
// Mutiple writers can be combined using MultiWriter.
//
// All writers have the same basic interface: a Write method,
// that can take any appropriately tagged struct; a Close method,
// which should be called to successfully complete writing; and
// a Cancel method, which should be called to abort writing and
// clean-up, in the event of an error or cancellation. It is safe
// to make mulltiple calls to Cancel, and it is safe to call Close
// after having previously called Cancel.
//
// All writers output their files atomically — that is to say:
// all output is written to a temporary location and only moved to the
// final output location when Close is called, meaning the output
// folder never contains any partially written files.
//
// Struct Tagging
//
// Structs to be used with peanut must have appropriately tagged
// fields, for example:
//  type Shape struct {
//  	ShapeID  string `peanut:"shape_id"`
//  	Name     string `peanut:"name"`
//  	NumSides int    `peanut:"num_sides"`
//  }
//
//  type Color struct {
//  	ColorID  string `peanut:"color_id"`
//  	Name     string `peanut:"name"`
//  	RBG      string `peanut:"rgb"`
//  }
//
// Fields without tags do not get written as output.
//
// Usage
//
// First create a writer, for example:
//  w := peanut.NewCSVWriter("/some/path/my-", "-data")
//
// Next, write some records to it:
// 	x := &Shape{
// 		ShapeID:  "sid1",
// 		Name:     "Square",
// 		NumSides: 4,
// 	}
// 	err := w.Write(x)
// 	// ...
//
// 	y := &Color{
// 		ColorID: "cid1",
// 		Name:    "red",
// 		RGB:     "ff0000",
// 	}
// 	err = w.Write(y)
// 	// ...
//
// 	z := &Shape{
// 		ShapeID:  "sid2",
// 		Name:     "Octogon",
// 		NumSides: 8,
// 	}
// 	err = w.Write(z)
// 	// ...
//
// When successfully completed:
//  err = w.Close()
//
//  // Output files will be:
//  // /some/path/my-Shape-data.csv
//  // /some/path/my-Color-data.csv
//
// Or, to abort the whole operation in the event of an error or cancellation while writing records:
//  err = w.Cancel()
//
// MultiWriter
//
// Multiple writers can be combined using MultiWriter:
//  w1 := peanut.NewCSVWriter("/some/path/my-", "-data")
//  w2 := peanut.NewExcelWriter("/some/path/my-", "-data")
//  w3 := &peanut.LogWriter{}
//  w := peanut.MultiWriter(w1, w2, w3)
// Here w will write records to CSV files, Excel files, and a logger.
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
//
// Tagging a field that has an unsupported datatype will result in a
// error when Write is called.
package peanut
