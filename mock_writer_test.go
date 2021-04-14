package peanut_test

import (
	"github.com/jimsmart/peanut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockWriter", func() {

	It("should capture the correct headers and data when structs are written", func() {

		expectedHeaders := map[string][]string{
			"Foo": {"foo_string1", "foo_int1"},
			"Bar": {"bar_int2", "bar_string2"},
		}

		expectedData := map[string][]map[string]string{
			"Foo": {
				{"foo_string1": "test 1", "foo_int1": "1"},
				{"foo_string1": "test 2", "foo_int1": "2"},
				{"foo_string1": "test 3", "foo_int1": "3"},
			},
			"Bar": {
				{"bar_int2": "1", "bar_string2": "test 1"},
				{"bar_int2": "2", "bar_string2": "test 2"},
				{"bar_int2": "3", "bar_string2": "test 3"},
			},
		}

		var err error
		w := &peanut.MockWriter{}
		for i := range testOutputFoo {
			err = w.Write(testOutputFoo[i])
			Expect(err).To(BeNil())
		}
		for i := range testOutputBar {
			err = w.Write(testOutputBar[i])
			Expect(err).To(BeNil())
		}
		err = w.Close()
		Expect(err).To(BeNil())

		Expect(w.Headers).To(Equal(expectedHeaders))
		Expect(w.Data).To(Equal(expectedData))

		Expect(w.CalledWrite).To(Equal(6))
		Expect(w.CalledClose).To(Equal(1))
		Expect(w.CalledCancel).To(Equal(0))

		// A bit lazy, arguably counters should have their own test - but MockWriter is simple.
		w.Cancel()
		Expect(w.CalledCancel).To(Equal(1))
	})
})
