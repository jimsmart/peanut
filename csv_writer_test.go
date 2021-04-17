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

	expectedOutput1 := "foo_string1,foo_int1\n" +
		"test 1,1\n" +
		"test 2,2\n" +
		"test 3,3\n"

	expectedOutput2 := "bar_int2,bar_string2\n" +
		"1,test 1\n" +
		"2,test 2\n" +
		"3,test 3\n"

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.csv")
		os.Remove("./test/output-Bar-sequential.csv")
		os.Remove("./test/output-Foo-interleave.csv")
		os.Remove("./test/output-Bar-interleave.csv")
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
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.csv").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.csv").ToNot(BeAnExistingFile())
	})

})
