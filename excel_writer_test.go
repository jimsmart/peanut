package peanut_test

import (
	"os"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
)

var _ = Describe("XecelWriter", func() {

	newFn := func(suffix string) peanut.Writer {
		w := peanut.NewExcelWriter("./test/output-", suffix)
		return w
	}

	expectedOutput1 := [][]string{
		{"foo_string1", "foo_int1"},
		{"test 1", "1"},
		{"test 2", "2"},
		{"test 3", "3"},
	}

	expectedOutput2 := [][]string{
		{"bar_int2", "bar_string2"},
		{"1", "test 1"},
		{"2", "test 2"},
		{"3", "test 3"},
	}

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.xlsx")
		os.Remove("./test/output-Bar-sequential.xlsx")
		os.Remove("./test/output-Foo-interleave.xlsx")
		os.Remove("./test/output-Bar-interleave.xlsx")
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
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.xlsx").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.xlsx").ToNot(BeAnExistingFile())
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
