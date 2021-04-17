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

	expectedOutput1 := `{"foo_int1":1,"foo_string1":"test 1"}` + "\n" +
		`{"foo_int1":2,"foo_string1":"test 2"}` + "\n" +
		`{"foo_int1":3,"foo_string1":"test 3"}` + "\n"

	expectedOutput2 := `{"bar_int2":1,"bar_string2":"test 1"}` + "\n" +
		`{"bar_int2":2,"bar_string2":"test 2"}` + "\n" +
		`{"bar_int2":3,"bar_string2":"test 3"}` + "\n"

	AfterEach(func() {
		os.Remove("./test/output-Foo-sequential.jsonl")
		os.Remove("./test/output-Bar-sequential.jsonl")
		os.Remove("./test/output-Foo-interleave.jsonl")
		os.Remove("./test/output-Bar-interleave.jsonl")
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
	})

	It("should not write anything when structs are written and cancel is called", func() {
		w := newFn("-cancel")

		testWritesAndCancel(w)

		Expect("./test/output-Foo-cancel.jsonl").ToNot(BeAnExistingFile())
		Expect("./test/output-Bar-cancel.jsonl").ToNot(BeAnExistingFile())
	})

})
