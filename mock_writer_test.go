package peanut_test

import (
	"github.com/jimsmart/peanut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockWriter", func() {

	It("should capture the correct headers and data when structs are written", func() {

		type Foo struct {
			StringField string `peanut:"test_string1,pk"`
			IntField    int    `peanut:"test_int1"`
		}

		type Bar struct {
			IntField    int    `peanut:"test_int2"`
			StringField string `peanut:"test_string2"`
		}

		testOutput1 := []*Foo{
			{StringField: "test 1", IntField: 1},
			{StringField: "test 2", IntField: 2},
			{StringField: "test 3", IntField: 3},
		}

		testOutput2 := []*Bar{
			{IntField: 1, StringField: "test 1"},
			{IntField: 2, StringField: "test 2"},
			{IntField: 3, StringField: "test 3"},
		}

		expectedHeaders := map[string][]string{
			"Foo": {"test_string1", "test_int1"},
			"Bar": {"test_int2", "test_string2"},
		}

		expectedData := map[string][]map[string]string{
			"Foo": {
				{"test_string1": "test 1", "test_int1": "1"},
				{"test_string1": "test 2", "test_int1": "2"},
				{"test_string1": "test 3", "test_int1": "3"},
			},
			"Bar": {
				{"test_int2": "1", "test_string2": "test 1"},
				{"test_int2": "2", "test_string2": "test 2"},
				{"test_int2": "3", "test_string2": "test 3"},
			},
		}

		var err error
		w := &peanut.MockWriter{}
		for i := range testOutput1 {
			err = w.Write(testOutput1[i])
			Expect(err).To(BeNil())
		}
		for i := range testOutput2 {
			err = w.Write(testOutput2[i])
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
