package peanut_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
)

var _ = Describe("TSVWriter", func() {

	newFn := func(suffix string) peanut.Writer {
		w := peanut.NewTSVWriter("./test/output-", suffix)
		return w
	}

	expectedOutput1 := "foo_string\tfoo_int\n" +
		"test 1\t1\n" +
		"test 2\t2\n" +
		"test 3\t3\n"

	expectedOutput2 := "bar_int\tbar_string\n" +
		"1\ttest 1\n" +
		"2\ttest 2\n" +
		"3\ttest 3\n"

	expectedOutput3 := "baz_string\tbaz_bool\tbaz_float32\tbaz_float64\tbaz_int\tbaz_int8\tbaz_int16\tbaz_int32\tbaz_int64\tbaz_uint\tbaz_uint8\tbaz_uint16\tbaz_uint32\tbaz_uint64\n" +
		"test 1\ttrue\t1.234\t9.876\t-12345\t-8\t-16\t-32\t-64\t12345\t8\t16\t32\t64\n"

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.tsv")
		os.Remove("./test/output-Bar-sequential.tsv")
		os.Remove("./test/output-Baz-sequential.tsv")
		os.Remove("./test/output-Foo-interleave.tsv")
		os.Remove("./test/output-Bar-interleave.tsv")
		os.Remove("./test/output-Baz-interleave.tsv")
	})

	It("should write the correct data when sequential structs are written", func() {
		w := newFn("-sequential")

		testWritesAndCloseSequential(w)

		output1, err := ioutil.ReadFile("./test/output-Foo-sequential.tsv")
		Expect(err).To(BeNil())
		Expect(string(output1)).To(Equal(expectedOutput1))

		output2, err := ioutil.ReadFile("./test/output-Bar-sequential.tsv")
		Expect(err).To(BeNil())
		Expect(string(output2)).To(Equal(expectedOutput2))

		output3, err := ioutil.ReadFile("./test/output-Baz-sequential.tsv")
		Expect(err).To(BeNil())
		Expect(string(output3)).To(Equal(expectedOutput3))
	})

	It("should write the correct data when interleaved structs are written", func() {
		w := newFn("-interleave")

		testWritesAndCloseSequential(w)

		output1, err := ioutil.ReadFile("./test/output-Foo-interleave.tsv")
		Expect(err).To(BeNil())
		Expect(string(output1)).To(Equal(expectedOutput1))

		output2, err := ioutil.ReadFile("./test/output-Bar-interleave.tsv")
		Expect(err).To(BeNil())
		Expect(string(output2)).To(Equal(expectedOutput2))

		output3, err := ioutil.ReadFile("./test/output-Baz-interleave.tsv")
		Expect(err).To(BeNil())
		Expect(string(output3)).To(Equal(expectedOutput3))
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.tsv").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.tsv").ToNot(BeAnExistingFile())
	})

	It("should return an error when the path is bad", func() {
		w := peanut.NewTSVWriter("./no-such-location/output-bogus-", "")

		err := w.Write(testOutputFoo[0])
		Expect(err).To(BeNil())

		err = w.Close()
		Expect(err).ToNot(BeNil())
	})

	Context("when given a struct with an unsupported field type", func() {

		It("should return an error with an informative message", func() {
			w := peanut.NewTSVWriter("./no-such-location/output-bogus-", "")
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
