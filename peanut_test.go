package peanut_test

type Foo struct {
	StringField  string `peanut:"foo_string1"`
	IntField     int    `peanut:"foo_int1"`
	IgnoredField int
}

type Bar struct {
	IntField    int    `peanut:"bar_int2"`
	StringField string `peanut:"bar_string2"`
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
