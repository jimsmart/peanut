package peanut_test

import (
	. "github.com/onsi/ginkgo"

	"github.com/jimsmart/peanut"
)

var _ = Describe("DiscardWriter", func() {

	It("should not return any errors when sequential structs are written", func() {
		w := &peanut.DiscardWriter{}
		testWritesAndCloseSequential(w)
	})

	It("should not return any errors when structs are written and cancel is called", func() {
		w := &peanut.DiscardWriter{}
		testWritesAndCancel(w)
	})
})
