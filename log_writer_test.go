package peanut_test

import (
	"bytes"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
)

var _ = Describe("LogWriter", func() {

	newFn := func() (peanut.Writer, *bytes.Buffer) {
		buf := &bytes.Buffer{}
		w := &peanut.LogWriter{
			Logger:  log.New(buf, "", 0),
			Verbose: true,
		}
		return w, buf
	}

	It("should write the correct text when sequential structs are written", func() {
		w, buf := newFn()

		testWritesAndCloseSequential(w)

		expectedOutput := "<Foo> foo_string: test 1 foo_int: 1\n" +
			"<Foo> foo_string: test 2 foo_int: 2\n" +
			"<Foo> foo_string: test 3 foo_int: 3\n" +
			"<Bar> bar_int: 1 bar_string: test 1\n" +
			"<Bar> bar_int: 2 bar_string: test 2\n" +
			"<Bar> bar_int: 3 bar_string: test 3\n" +
			"<Baz> baz_string: test 1 baz_bool: true baz_float32: 1.234 baz_float64: 9.876 baz_int: -12345 baz_int8: -8 baz_int16: -16 baz_int32: -32 baz_int64: -64 baz_uint: 12345 baz_uint8: 8 baz_uint16: 16 baz_uint32: 32 baz_uint64: 64\n" +
			"Called LogWriter.Close\n"

		output := string(buf.Bytes())
		// fmt.Println(output)
		Expect(output).To(Equal(expectedOutput))
	})

	It("should write the correct text when interleaved structs are written", func() {
		w, buf := newFn()

		testWritesAndCloseInterleaved(w)

		expectedOutput := "<Foo> foo_string: test 1 foo_int: 1\n" +
			"<Bar> bar_int: 1 bar_string: test 1\n" +
			"<Foo> foo_string: test 2 foo_int: 2\n" +
			"<Bar> bar_int: 2 bar_string: test 2\n" +
			"<Foo> foo_string: test 3 foo_int: 3\n" +
			"<Bar> bar_int: 3 bar_string: test 3\n" +
			"<Baz> baz_string: test 1 baz_bool: true baz_float32: 1.234 baz_float64: 9.876 baz_int: -12345 baz_int8: -8 baz_int16: -16 baz_int32: -32 baz_int64: -64 baz_uint: 12345 baz_uint8: 8 baz_uint16: 16 baz_uint32: 32 baz_uint64: 64\n" +
			"Called LogWriter.Close\n"

		output := string(buf.Bytes())
		// fmt.Println(output)
		Expect(output).To(Equal(expectedOutput))
	})

	It("should write the correct text when structs are written and cancel is called", func() {
		w, buf := newFn()

		testWritesAndCancel(w)

		expectedOutput := "<Foo> foo_string: test 1 foo_int: 1\n" +
			"<Bar> bar_int: 1 bar_string: test 1\n" +
			"Called LogWriter.Cancel\n"

		output := string(buf.Bytes())
		Expect(output).To(Equal(expectedOutput))
	})

})
