package peanut_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
)

var _ = Describe("JSONLWriter", func() {

	newFn := func(suffix string) peanut.Writer {
		w := peanut.NewJSONLWriter("./test/output-", suffix)
		return w
	}

	expectedOutput1 := `{"foo_int":1,"foo_string":"test 1"}` + "\n" +
		`{"foo_int":2,"foo_string":"test 2"}` + "\n" +
		`{"foo_int":3,"foo_string":"test 3"}` + "\n"

	expectedOutput2 := `{"bar_int":1,"bar_string":"test 1"}` + "\n" +
		`{"bar_int":2,"bar_string":"test 2"}` + "\n" +
		`{"bar_int":3,"bar_string":"test 3"}` + "\n"

	expectedOutput3 := `{"baz_bool":true,"baz_float32":1.234,"baz_float64":9.876,"baz_int":-12345,"baz_int16":-16,"baz_int32":-32,"baz_int64":-64,"baz_int8":-8,"baz_string":"test 1","baz_uint":12345,"baz_uint16":16,"baz_uint32":32,"baz_uint64":64,"baz_uint8":8}` + "\n"

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.jsonl")
		os.Remove("./test/output-Bar-sequential.jsonl")
		os.Remove("./test/output-Baz-sequential.jsonl")
		os.Remove("./test/output-Foo-interleave.jsonl")
		os.Remove("./test/output-Bar-interleave.jsonl")
		os.Remove("./test/output-Baz-interleave.jsonl")
	})

	It("should write the correct data when sequential structs are written", func() {
		w := newFn("-sequential")

		testWritesAndCloseSequential(w)

		output1, err := ioutil.ReadFile("./test/output-Foo-sequential.jsonl")
		Expect(err).To(BeNil())
		Expect(string(output1)).To(Equal(expectedOutput1))

		output2, err := ioutil.ReadFile("./test/output-Bar-sequential.jsonl")
		Expect(err).To(BeNil())
		Expect(string(output2)).To(Equal(expectedOutput2))

		output3, err := ioutil.ReadFile("./test/output-Baz-sequential.jsonl")
		Expect(err).To(BeNil())
		// fmt.Println(string(output3))
		Expect(string(output3)).To(Equal(expectedOutput3))
	})

	It("should write the correct data when interleaved structs are written", func() {
		w := newFn("-interleave")

		testWritesAndCloseSequential(w)

		output1, err := ioutil.ReadFile("./test/output-Foo-interleave.jsonl")
		Expect(err).To(BeNil())
		Expect(string(output1)).To(Equal(expectedOutput1))

		output2, err := ioutil.ReadFile("./test/output-Bar-interleave.jsonl")
		Expect(err).To(BeNil())
		Expect(string(output2)).To(Equal(expectedOutput2))

		output3, err := ioutil.ReadFile("./test/output-Baz-interleave.jsonl")
		Expect(err).To(BeNil())
		// fmt.Println(string(output3))
		Expect(string(output3)).To(Equal(expectedOutput3))
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.jsonl").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.jsonl").ToNot(BeAnExistingFile())
	})

	It("should return an error when Write is called after Close", func() {
		w := newFn("-close-write")

		testWriteAfterClose(w)

		Expect("./test/output-Foo-close-write.jsonl").ToNot(BeAnExistingFile())
	})

	It("should return an error when the path is bad", func() {
		w := peanut.NewJSONLWriter("./no-such-location/output-bogus-", "")

		err := w.Write(testOutputFoo[0])
		Expect(err).To(BeNil())

		err = w.Close()
		Expect(err).ToNot(BeNil())
	})

	Context("when given a struct with an unsupported field type", func() {

		It("should return an error with an informative message", func() {
			w := peanut.NewJSONLWriter("./no-such-location/output-bogus-", "")
			defer func() {
				err1 := w.Cancel()
				err2 := w.Close()
				Expect(err1).To(BeNil())
				Expect(err2).To(BeNil())
			}()

			err := w.Write(BadUnsupported{})
			Expect(err).ToNot(BeNil())

			// Expect error message to be informative.
			Expect(err.Error()).To(SatisfyAll(
				MatchRegexp(`slice`),          // type
				MatchRegexp("BytesField"),     // field name
				MatchRegexp("BadUnsupported"), // struct name
			))
		})
	})
})
