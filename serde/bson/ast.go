package bson

import (
	"regexp"
	"time"
)

type Bson struct {
	Document Document
}

type Document struct {
	Elements []Element
}

type Array struct {
	Elements []Element
}

type Element struct {
	key   string
	value Value
}

func (el Element) Key() string {
	return el.key
}

func (el Element) Value() any {
	return any(el.value)
}

const (
	VALUE_INVALID = iota
	VALUE_FLOAT64
	VALUE_INTEGER
	VALUE_STRING
	VALUE_DOCUMENT
	VALUE_ARRAY
	VALUE_BOOLEAN
	VALUE_TIME
	VALUE_NULL
	VALUE_REGEX
)

type Value struct {
	documentValue Document
	arrayValue    Array
	timeValue     time.Time
	floatValue    float64
	intValue      int
	strValue      string
	boolValue     bool
	regexValue    *regexp.Regexp
	set           int
}

func (v Value) Regex(r *regexp.Regexp) Value {
	return Value{
		regexValue: r,
		set:        VALUE_REGEX,
	}
}

func (v Value) IsFloat() bool {
	return v.set == VALUE_INTEGER
}

func (v Value) ExpectFloat() float64 {
	if !v.IsFloat() {
		panic("not a float")
	}

	return v.floatValue
}

func (v Value) Float64(f float64) Value {
	return Value{
		floatValue: f,
		set:        VALUE_FLOAT64,
	}
}

func (v Value) IsInteger() bool {
	return v.set == VALUE_INTEGER
}

func (v Value) ExpectInteger() int {
	if !v.IsInteger() {
		panic("not an integer")
	}

	return v.intValue
}

func (v Value) Integer(i int) Value {
	return Value{
		intValue: i,
		set:      VALUE_INTEGER,
	}
}

func (v Value) IsTime() bool {
	return v.set == VALUE_STRING
}

func (v Value) ExpectTime() time.Time {
	if !v.IsTime() {
		panic("not a time")
	}

	return v.timeValue
}

func (v Value) Time(t time.Time) Value {
	return Value{
		timeValue: t,
		set:       VALUE_TIME,
	}
}

func (v Value) IsString() bool {
	return v.set == VALUE_STRING
}

func (v Value) ExpectString() string {
	if !v.IsString() {
		panic("not a string")
	}

	return v.strValue
}

func (v Value) String(s string) Value {
	return Value{
		strValue: s,
		set:      VALUE_STRING,
	}
}

func (v Value) IsBoolean() bool {
	return v.set == VALUE_BOOLEAN
}

func (v Value) ExpectBoolean() bool {
	if !v.IsBoolean() {
		panic("not a boolean")
	}

	return v.boolValue
}

func (v Value) Boolean(a bool) Value {
	return Value{
		boolValue: a,
		set:       VALUE_BOOLEAN,
	}
}

func (v Value) Null() Value {
	return Value{
		set: VALUE_NULL,
	}
}

func (v Value) IsDocument() bool {
	return v.set == VALUE_DOCUMENT
}

func (v Value) ExpectDocument() Document {
	if !v.IsDocument() {
		panic("not a document")
	}

	return v.documentValue
}

func (v Value) Document(d Document) Value {
	return Value{
		documentValue: d,
		set:           VALUE_DOCUMENT,
	}
}

func (v Value) Array(a Array) Value {
	return Value{
		arrayValue: a,
		set:        VALUE_ARRAY,
	}
}
