package peanut_test

import (
	"github.com/jimsmart/peanut"
	. "github.com/onsi/gomega"
)

type Foo struct {
	StringField  string `peanut:"foo_string1,pk"`
	IntField     int    `peanut:"foo_int1"`
	IgnoredField int
}

type Bar struct {
	IntField    int    `peanut:"bar_int2,pk"`
	StringField string `peanut:"bar_string2,pk"`
}

var testOutputFoo = []*Foo{
	{StringField: "test 1", IntField: 1},
	{StringField: "test 2", IntField: 2},
	{StringField: "test 3", IntField: 3},
}

var testOutputBar = []*Bar{
	{IntField: 1, StringField: "test 1"},
	{IntField: 2, StringField: "test 2"},
	{IntField: 3, StringField: "test 3"},
}

func testWritesAndCloseSequential(w peanut.Writer) {
	var err error
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
}

func testWritesAndCloseInterleaved(w peanut.Writer) {
	var err error
	for i := range testOutputFoo {
		err = w.Write(testOutputFoo[i])
		Expect(err).To(BeNil())
		err = w.Write(testOutputBar[i])
		Expect(err).To(BeNil())
	}
	err = w.Close()
	Expect(err).To(BeNil())
}

func testWritesAndCancel(w peanut.Writer) {
	var err error
	err = w.Write(testOutputFoo[0])
	Expect(err).To(BeNil())
	err = w.Write(testOutputBar[0])
	Expect(err).To(BeNil())
	err = w.Cancel()
	Expect(err).To(BeNil())
}
