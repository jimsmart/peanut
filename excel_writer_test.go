package peanut_test

import (
	"os"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
)

var _ = Describe("ExcelWriter", func() {

	newFn := func(suffix string) peanut.Writer {
		w := peanut.NewExcelWriter("./test/output-", suffix)
		return w
	}

	expectedOutput1 := [][]string{
		{"foo_string", "foo_int"},
		{"test 1", "1"},
		{"test 2", "2"},
		{"test 3", "3"},
	}

	expectedOutput2 := [][]string{
		{"bar_int", "bar_string"},
		{"1", "test 1"},
		{"2", "test 2"},
		{"3", "test 3"},
	}

	expectedOutput3 := [][]string{
		{"baz_string", "baz_bool", "baz_float32", "baz_float64", "baz_int", "baz_int8", "baz_int16", "baz_int32", "baz_int64", "baz_uint", "baz_uint8", "baz_uint16", "baz_uint32", "baz_uint64"},
		{"test 1", "1", "1.234", "9.876", "-12345", "-8", "-16", "-32", "-64", "12345", "8", "16", "32", "64"},
	}

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.xlsx")
		os.Remove("./test/output-Bar-sequential.xlsx")
		os.Remove("./test/output-Baz-sequential.xlsx")
		os.Remove("./test/output-Foo-interleave.xlsx")
		os.Remove("./test/output-Bar-interleave.xlsx")
		os.Remove("./test/output-Baz-interleave.xlsx")
	})

	It("should write the correct data when sequential structs are written", func() {
		w := newFn("-sequential")

		testWritesAndCloseSequential(w)

		output1, err := readExcel("./test/output-Foo-sequential.xlsx")
		Expect(err).To(BeNil())
		Expect(output1).To(Equal(expectedOutput1))

		output2, err := readExcel("./test/output-Bar-sequential.xlsx")
		Expect(err).To(BeNil())
		Expect(output2).To(Equal(expectedOutput2))

		output3, err := readExcel("./test/output-Baz-sequential.xlsx")
		Expect(err).To(BeNil())
		Expect(output3).To(Equal(expectedOutput3))
	})

	It("should write the correct data when interleaved structs are written", func() {
		w := newFn("-interleave")

		testWritesAndCloseSequential(w)

		output1, err := readExcel("./test/output-Foo-interleave.xlsx")
		Expect(err).To(BeNil())
		Expect(output1).To(Equal(expectedOutput1))

		output2, err := readExcel("./test/output-Bar-interleave.xlsx")
		Expect(err).To(BeNil())
		Expect(output2).To(Equal(expectedOutput2))

		output3, err := readExcel("./test/output-Baz-interleave.xlsx")
		Expect(err).To(BeNil())
		Expect(output3).To(Equal(expectedOutput3))
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.xlsx").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.xlsx").ToNot(BeAnExistingFile())
	})

	It("should return an error when Write is called after Close", func() {
		w := newFn("-close-write")

		testWriteAfterClose(w)

		Expect("./test/output-Foo-close-write.xlsx").ToNot(BeAnExistingFile())
	})

	It("should return an error when the path is bad", func() {
		w := peanut.NewExcelWriter("./no-such-location/output-bogus-", "")

		err := w.Write(testOutputFoo[0])
		Expect(err).To(BeNil())

		err = w.Close()
		Expect(err).ToNot(BeNil())
	})

	Context("when given a struct with an unsupported field type", func() {

		It("should return an error with an informative message", func() {
			w := peanut.NewExcelWriter("./no-such-location/output-bogus-", "")
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

func readExcel(filename string) ([][]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	var out [][]string
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		out = append(out, row)
	}
	return out, nil
}
