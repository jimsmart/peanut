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

		expectedOutput := "<Foo> foo_string1: test 1 foo_int1: 1\n" +
			"<Foo> foo_string1: test 2 foo_int1: 2\n" +
			"<Foo> foo_string1: test 3 foo_int1: 3\n" +
			"<Bar> bar_int2: 1 bar_string2: test 1\n" +
			"<Bar> bar_int2: 2 bar_string2: test 2\n" +
			"<Bar> bar_int2: 3 bar_string2: test 3\n" +
			"Called LogWriter.Close\n"

		output := string(buf.Bytes())
		Expect(output).To(Equal(expectedOutput))
	})

	It("should write the correct text when interleaved structs are written", func() {
		w, buf := newFn()

		testWritesAndCloseInterleaved(w)

		expectedOutput := "<Foo> foo_string1: test 1 foo_int1: 1\n" +
			"<Bar> bar_int2: 1 bar_string2: test 1\n" +
			"<Foo> foo_string1: test 2 foo_int1: 2\n" +
			"<Bar> bar_int2: 2 bar_string2: test 2\n" +
			"<Foo> foo_string1: test 3 foo_int1: 3\n" +
			"<Bar> bar_int2: 3 bar_string2: test 3\n" +
			"Called LogWriter.Close\n"

		output := string(buf.Bytes())
		Expect(output).To(Equal(expectedOutput))
	})

	It("should write the correct text when structs are written and cancel is called", func() {
		w, buf := newFn()

		testWritesAndCancel(w)

		expectedOutput := "<Foo> foo_string1: test 1 foo_int1: 1\n" +
			"<Bar> bar_int2: 1 bar_string2: test 1\n" +
			"Called LogWriter.Cancel\n"

		output := string(buf.Bytes())
		Expect(output).To(Equal(expectedOutput))
	})

})
