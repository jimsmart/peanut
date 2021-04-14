package peanut_test

import (
	"bytes"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/elektra/scraper/peanut"
)

var _ = Describe("LogWriter", func() {

	It("should write the correct text when structs are written", func() {

		expectedOutput := "<Foo> foo_string1: test 1 foo_int1: 1\n" +
			"<Foo> foo_string1: test 2 foo_int1: 2\n" +
			"<Foo> foo_string1: test 3 foo_int1: 3\n" +
			"Called LogWriter.Close\n" +
			"Called LogWriter.Cancel\n"

		var err error

		buf := &bytes.Buffer{}

		w := &peanut.LogWriter{
			Logger:  log.New(buf, "", 0),
			Verbose: true,
		}

		for i := range testOutputFoo {
			err = w.Write(testOutputFoo[i])
			Expect(err).To(BeNil())
		}
		err = w.Close()
		Expect(err).To(BeNil())

		// This only logs a message for LogWriter.
		// This is not actually proper usage of the API in any way.
		w.Cancel()

		output := string(buf.Bytes())
		Expect(output).To(Equal(expectedOutput))
	})

})
