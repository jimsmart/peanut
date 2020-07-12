package peanut_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPeanut(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "peanut Suite")
}
