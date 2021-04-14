package peanut_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/elektra/scraper/peanut"
)

var _ = Describe("MultiWriter", func() {

	It("should call the methods of the writers it wraps", func() {

		testOutput := []*Foo{
			{StringField: "test 1", IntField: 1},
			{StringField: "test 2", IntField: 2},
			{StringField: "test 3", IntField: 3},
		}

		var err error
		w1 := &peanut.MockWriter{}
		w2 := &peanut.MockWriter{}
		w := peanut.MultiWriter(w1, w2)
		for i := range testOutput {
			err = w.Write(testOutput[i])
			Expect(err).To(BeNil())
		}
		err = w.Close()
		Expect(err).To(BeNil())

		// This is a no-op on MockWriter, and just increments a counter.
		//  This is not meant to illustrate correct usage of this API.
		w.Cancel()

		Expect(w1.CalledWrite).To(Equal(3))
		Expect(w1.CalledClose).To(Equal(1))
		Expect(w1.CalledCancel).To(Equal(1))

		Expect(w1).To(Equal(w2))
	})
})
