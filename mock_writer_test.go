package peanut_test

import (
	"github.com/jimsmart/peanut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockWriter", func() {

	expectedHeaders := map[string][]string{
		"Foo": {"foo_string", "foo_int"},
		"Bar": {"bar_int", "bar_string"},
		"Baz": {"baz_string", "baz_bool", "baz_float32", "baz_float64", "baz_int", "baz_int8", "baz_int16", "baz_int32", "baz_int64", "baz_uint", "baz_uint8", "baz_uint16", "baz_uint32", "baz_uint64"},
	}

	expectedData := map[string][]map[string]string{
		"Foo": {
			{"foo_string": "test 1", "foo_int": "1"},
			{"foo_string": "test 2", "foo_int": "2"},
			{"foo_string": "test 3", "foo_int": "3"},
		},
		"Bar": {
			{"bar_int": "1", "bar_string": "test 1"},
			{"bar_int": "2", "bar_string": "test 2"},
			{"bar_int": "3", "bar_string": "test 3"},
		},
		"Baz": {
			{
				"baz_string":  "test 1",
				"baz_bool":    "true",
				"baz_float32": "1.234",
				"baz_float64": "9.876",
				"baz_int":     "-12345",
				"baz_int8":    "-8",
				"baz_int16":   "-16",
				"baz_int32":   "-32",
				"baz_int64":   "-64",
				"baz_uint":    "12345",
				"baz_uint8":   "8",
				"baz_uint16":  "16",
				"baz_uint32":  "32",
				"baz_uint64":  "64",
			},
		},
	}

	expectedDataPartial := map[string][]map[string]string{
		"Bar": {
			{"bar_int": "1", "bar_string": "test 1"},
			{"bar_int": "2", "bar_string": "test 2"},
			{"bar_int": "3", "bar_string": "test 3"},
		},
	}

	It("should capture the correct headers and data when sequential structs are written", func() {
		w := &peanut.MockWriter{}

		testWritesAndCloseSequential(w)

		Expect(w.Headers).To(Equal(expectedHeaders))
		Expect(w.Data).To(Equal(expectedData))

		Expect(w.CalledWrite).To(Equal(10))
		Expect(w.CalledClose).To(Equal(1))
		Expect(w.CalledCancel).To(Equal(1))
	})

	It("should capture the correct headers and data when interleaved structs are written", func() {
		w := &peanut.MockWriter{}

		testWritesAndCloseInterleaved(w)

		Expect(w.Headers).To(Equal(expectedHeaders))
		Expect(w.Data).To(Equal(expectedData))

		Expect(w.CalledWrite).To(Equal(10))
		Expect(w.CalledClose).To(Equal(1))
		Expect(w.CalledCancel).To(Equal(0))
	})

	It("should have the correct counters when structs are written and cancel is called", func() {
		w := &peanut.MockWriter{}

		testWritesAndCancel(w)

		Expect(w.CalledCancel).To(Equal(1))
		Expect(w.CalledClose).To(Equal(1))
	})

	It("should capture the correct headers and data when DisableDataCapture is used", func() {
		w := &peanut.MockWriter{}
		w.DisableDataCapture = map[string]bool{
			"Foo": true,
			"Baz": true,
		}

		testWritesAndCloseSequential(w)

		Expect(w.Headers).To(Equal(expectedHeaders))
		Expect(w.Data).To(Equal(expectedDataPartial))

		Expect(w.CalledWrite).To(Equal(10))
		Expect(w.CalledClose).To(Equal(1))
		Expect(w.CalledCancel).To(Equal(1))
	})

	Context("when given a struct with an unsupported field type", func() {

		It("should return an error with an informative message", func() {
			w := &peanut.MockWriter{}

			testWriteBadType(w)
			// TODO(js) Do we need further checks, e.g. file not exists ...?
		})
	})
})
