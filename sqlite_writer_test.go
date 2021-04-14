package peanut_test

import (
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
	// "github.com/jimsmart/peanut"
)

var _ = Describe("SQLiteWriter", func() {

	// XIt("should write the correct data when structs are written", func() {

	// 	type Foo struct {
	// 		StringField string `peanut:"test_string1,pk"`
	// 		IntField    int    `peanut:"test_int1"`
	// 		// IgnoredField string // TODO Currently all fields are processed. We should ignore those without appropriate tags.
	// 	}

	// 	testOutput := []*Foo{
	// 		{StringField: "test 1", IntField: 1},
	// 		{StringField: "test 2", IntField: 2},
	// 		{StringField: "test 3", IntField: 3},
	// 	}

	// 	// expectedOutput := "<Foo> test_string1: test 1 test_int1: 1\n" +
	// 	// 	"<Foo> test_string1: test 2 test_int1: 2\n" +
	// 	// 	"<Foo> test_string1: test 3 test_int1: 3\n" +
	// 	// 	"Called LogWriter.Close\n" +
	// 	// 	"Called LogWriter.Cancel\n"

	// 	var err error

	// 	w := peanut.NewSQLiteWriter("foo")

	// 	for i := range testOutput {
	// 		err = w.Write(testOutput[i])
	// 		Expect(err).To(BeNil())
	// 	}
	// 	err = w.Close()
	// 	Expect(err).To(BeNil())

	// 	// // This only logs a message for LogWriter.
	// 	// // This is not actually proper usage of the API in any way.
	// 	// w.Cancel()

	// 	// output := string(buf.Bytes())
	// 	// Expect(output).To(Equal(expectedOutput))
	// })

})
