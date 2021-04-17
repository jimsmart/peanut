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

	expectedOutput1 := "foo_string1\tfoo_int1\n" +
		"test 1\t1\n" +
		"test 2\t2\n" +
		"test 3\t3\n"

	expectedOutput2 := "bar_int2\tbar_string2\n" +
		"1\ttest 1\n" +
		"2\ttest 2\n" +
		"3\ttest 3\n"

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.tsv")
		os.Remove("./test/output-Bar-sequential.tsv")
		os.Remove("./test/output-Foo-interleave.tsv")
		os.Remove("./test/output-Bar-interleave.tsv")
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
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.tsv").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.tsv").ToNot(BeAnExistingFile())
	})

})
