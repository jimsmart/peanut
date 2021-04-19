package peanut_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jimsmart/peanut"
)

var _ = Describe("MultiWriter", func() {

	It("should call the methods of the writers it wraps", func() {

		w1 := &peanut.MockWriter{}
		w2 := &peanut.MockWriter{}
		w := peanut.MultiWriter(w1, w2)

		testWritesAndCloseSequential(w)

		Expect(w1.CalledWrite).To(Equal(10))
		Expect(w1.CalledClose).To(Equal(1))
		Expect(w1.CalledCancel).To(Equal(1))

		// Did it actually write to both?
		Expect(w1).To(Equal(w2))
	})
})
