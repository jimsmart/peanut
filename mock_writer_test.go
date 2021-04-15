package peanut_test

import (
	"github.com/jimsmart/peanut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockWriter", func() {

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

	It("should capture the correct headers and data when sequential structs are written", func() {
		w := &peanut.MockWriter{}

		testWritesAndCloseSequential(w)

		Expect(w.Headers).To(Equal(expectedHeaders))
		Expect(w.Data).To(Equal(expectedData))

		Expect(w.CalledWrite).To(Equal(6))
		Expect(w.CalledClose).To(Equal(1))
		Expect(w.CalledCancel).To(Equal(0))
	})

	It("should capture the correct headers and data when interleaved structs are written", func() {
		w := &peanut.MockWriter{}

		testWritesAndCloseInterleaved(w)

		Expect(w.Headers).To(Equal(expectedHeaders))
		Expect(w.Data).To(Equal(expectedData))

		Expect(w.CalledWrite).To(Equal(6))
		Expect(w.CalledClose).To(Equal(1))
		Expect(w.CalledCancel).To(Equal(0))
	})

	It("should have the correct counters when structs are written and cancel is called", func() {
		w := &peanut.MockWriter{}

		testWritesAndCancel(w)

		Expect(w.CalledClose).To(Equal(0))
		Expect(w.CalledCancel).To(Equal(1))
	})

})
