package peanut_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
)

var _ = Describe("CSVWriter", func() {

	newFn := func(suffix string) peanut.Writer {
		w := peanut.NewCSVWriter("./test/output-", suffix)
		return w
	}

	expectedOutput1 := "foo_string,foo_int\n" +
		"test 1,1\n" +
		"test 2,2\n" +
		"test 3,3\n"

	expectedOutput2 := "bar_int,bar_string\n" +
		"1,test 1\n" +
		"2,test 2\n" +
		"3,test 3\n"

	expectedOutput3 := "baz_string,baz_bool,baz_float32,baz_float64,baz_int,baz_int8,baz_int16,baz_int32,baz_int64,baz_uint,baz_uint8,baz_uint16,baz_uint32,baz_uint64\n" +
		"test 1,true,1.234,9.876,-12345,-8,-16,-32,-64,12345,8,16,32,64\n"

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.csv")
		os.Remove("./test/output-Bar-sequential.csv")
		os.Remove("./test/output-Baz-sequential.csv")
		os.Remove("./test/output-Foo-interleave.csv")
		os.Remove("./test/output-Bar-interleave.csv")
		os.Remove("./test/output-Baz-interleave.csv")
	})

	It("should write the correct data when sequential structs are written", func() {
		w := newFn("-sequential")

		testWritesAndCloseSequential(w)

		output1, err := ioutil.ReadFile("./test/output-Foo-sequential.csv")
		Expect(err).To(BeNil())
		Expect(string(output1)).To(Equal(expectedOutput1))

		output2, err := ioutil.ReadFile("./test/output-Bar-sequential.csv")
		Expect(err).To(BeNil())
		Expect(string(output2)).To(Equal(expectedOutput2))

		output3, err := ioutil.ReadFile("./test/output-Baz-sequential.csv")
		Expect(err).To(BeNil())
		Expect(string(output3)).To(Equal(expectedOutput3))
	})

	It("should write the correct data when interleaved structs are written", func() {
		w := newFn("-interleave")

		testWritesAndCloseSequential(w)

		output1, err := ioutil.ReadFile("./test/output-Foo-interleave.csv")
		Expect(err).To(BeNil())
		Expect(string(output1)).To(Equal(expectedOutput1))

		output2, err := ioutil.ReadFile("./test/output-Bar-interleave.csv")
		Expect(err).To(BeNil())
		Expect(string(output2)).To(Equal(expectedOutput2))

		output3, err := ioutil.ReadFile("./test/output-Baz-interleave.csv")
		Expect(err).To(BeNil())
		Expect(string(output3)).To(Equal(expectedOutput3))
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.csv").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.csv").ToNot(BeAnExistingFile())
	})

	It("should return an error when the path is bad", func() {
		w := peanut.NewCSVWriter("./no-such-location/output-bogus-", "")

		err := w.Write(testOutputFoo[0])
		Expect(err).To(BeNil())

		err = w.Close()
		Expect(err).ToNot(BeNil())
	})

	Context("when given a struct with an unsupported field type", func() {

		It("should panic with an appropriate message", func() {
			w := peanut.NewCSVWriter("./no-such-location/output-bogus-", "")

			// err := w.Write(&BadField{})
			// Expect(err).To(BeNil())

			fn := func() {
				w.Write(BadUnsupported{})
			}

			// Expect panic message to be informative.
			Expect(fn).To(PanicWith(SatisfyAll(
				MatchRegexp(`slice`),          // type
				MatchRegexp("BytesField"),     // field name
				MatchRegexp("BadUnsupported"), // struct name
			)))

			// Expect(err).To(BeNil())

			// err = w.Close()
			// Expect(err).ToNot(BeNil())
		})
	})

})
