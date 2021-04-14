package peanut_test

import (
	"github.com/jimsmart/elektra/scraper/peanut"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MockWriter", func() {

	It("should capture the correct headers and data when structs are written", func() {

		// type Foo struct {
		// 	StringField string `peanut:"test_string1,pk"`
		// 	IntField    int    `peanut:"test_int1"`
		// }

		// type Bar struct {
		// 	IntField    int    `peanut:"test_int2"`
		// 	StringField string `peanut:"test_string2"`
		// }

		// testOutput1 := []*Foo{
		// 	{StringField: "test 1", IntField: 1},
		// 	{StringField: "test 2", IntField: 2},
		// 	{StringField: "test 3", IntField: 3},
		// }

		// testOutput2 := []*Bar{
		// 	{IntField: 1, StringField: "test 1"},
		// 	{IntField: 2, StringField: "test 2"},
		// }

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
