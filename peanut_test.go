package peanut_test

import (
	"github.com/jimsmart/peanut"
	. "github.com/onsi/gomega"
)

type Foo struct {
	StringField string `peanut:"foo_string,pk"`
	IntField    int    `peanut:"foo_int"`
}

type Bar struct {
	IntField    int    `peanut:"bar_int,pk"`
	StringField string `peanut:"bar_string,pk"`
}

type Baz struct {
	StringField  string  `peanut:"baz_string,pk"`
	BoolField    bool    `peanut:"baz_bool"`
	Float32Field float64 `peanut:"baz_float32"`
	Float64Field float64 `peanut:"baz_float64"`
	IntField     int     `peanut:"baz_int"`
	Int8Field    int8    `peanut:"baz_int8"`
	Int16Field   int16   `peanut:"baz_int16"`
	Int32Field   int32   `peanut:"baz_int32"`
	Int64Field   int64   `peanut:"baz_int64"`
	UintField    uint    `peanut:"baz_uint"`
	Uint8Field   uint8   `peanut:"baz_uint8"`
	Uint16Field  uint16  `peanut:"baz_uint16"`
	Uint32Field  uint32  `peanut:"baz_uint32"`
	Uint64Field  uint64  `peanut:"baz_uint64"`
	IgnoredField int     // No tag.
	privateField int     `peanut:"baz_private"`
}

type Qux struct {
	StringField string
	IntField    int
}

type BadUnsupported struct {
	BytesField []byte `peanut:"bytes_field"`
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

var testOutputBaz = []Baz{
	{
		StringField:  "test 1",
		BoolField:    true,
		Float32Field: 1.234,
		Float64Field: 9.876,
		IntField:     -12345,
		Int8Field:    -8,
		Int16Field:   -16,
		Int32Field:   -32,
		Int64Field:   -64,
		UintField:    12345,
		Uint8Field:   8,
		Uint16Field:  16,
		Uint32Field:  32,
		Uint64Field:  64,
	},
}

var testOutputQux = []*Qux{
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
	for i := range testOutputBaz {
		err = w.Write(testOutputBaz[i])
		Expect(err).To(BeNil())
	}
	for i := range testOutputQux {
		err = w.Write(testOutputQux[i])
		Expect(err).To(BeNil())
	}
	err = w.Close()
	Expect(err).To(BeNil())

	// Calling Cancel after Close should be a no-op.
	err = w.Cancel()
	Expect(err).To(BeNil())
}

func testWritesAndCloseInterleaved(w peanut.Writer) {
	var err error
	for i := range testOutputFoo {
		err = w.Write(testOutputFoo[i])
		Expect(err).To(BeNil())
		err = w.Write(testOutputBar[i])
		Expect(err).To(BeNil())
		err = w.Write(testOutputQux[i])
		Expect(err).To(BeNil())
	}
	for i := range testOutputBaz {
		err = w.Write(testOutputBaz[i])
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

	// Calling Close after Cancel should be a no-op.
	err = w.Close()
	Expect(err).To(BeNil())
}

func testWriteBadType(w peanut.Writer) {

	defer func() {
		err1 := w.Cancel()
		err2 := w.Close()
		Expect(err1).To(BeNil())
		Expect(err2).To(BeNil())
	}()

	err := w.Write(BadUnsupported{})
	Expect(err).ToNot(BeNil())

	// Expect error message to be informative.
	Expect(err.Error()).To(SatisfyAll(
		MatchRegexp(`slice`),          // type
		MatchRegexp("BytesField"),     // field name
		MatchRegexp("BadUnsupported"), // struct name
	))
}

func testWriteAfterClose(w peanut.Writer) {
	var err error
	err = w.Close()
	Expect(err).To(BeNil())
	err = w.Write(testOutputFoo[0])
	Expect(err).To(Equal(peanut.ErrClosedWriter))
}
