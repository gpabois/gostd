package json

type Json struct {
	array    Array
	document Document
	set      int
}

func (json Json) IsArray() bool {
	return json.set == 1
}

func (json Json) ExpectArray() Array {
	if !json.IsArray() {
		panic("not an array")
	}

	return json.array
}

func (json Json) Array(array Array) Json {
	json.array = array
	json.set = 1
	return json
}

func (json Json) Document(document Document) Json {
	json.document = document
	json.set = 2
	return json
}

func (json Json) IsDocument() bool {
	return json.set == 1
}

func (json Json) ExpectDocument() Document {
	if !json.IsDocument() {
		panic("not a document")
	}

	return json.document
}

type Document struct {
	Pairs []Element
}

type Element struct {
	key   string
	value Value
}

func (el Element) Key() string {
	return el.key
}

func (el Element) Value() any {
	return el.value
}

type Array struct {
	Elements []Value
}

const (
	VALUE_DOCUMENT = iota
	VALUE_ARRAY
	VALUE_BOOLEAN
	VALUE_STRING
	VALUE_INTEGER
	VALUE_FLOAT
	VALUE_NULL
)

type Value struct {
	documentValue Document
	arrayValue    Array
	boolValue     bool
	stringValue   string
	integerValue  int
	floatValue    float64

	// Define the value set (similarly to union)
	set int
}

func (val Value) IsDocument() bool {
	return val.set == VALUE_DOCUMENT
}

func (val Value) ExpectDocument() Document {
	if !val.IsDocument() {
		panic("not a document")
	}

	return val.documentValue
}

func (val Value) Document(document Document) Value {
	return Value{
		documentValue: document,
		set:           VALUE_DOCUMENT,
	}
}

func (val Value) IsArray() bool {
	return val.set == VALUE_ARRAY
}

func (val Value) ExpectArray() Array {
	if !val.IsArray() {
		panic("not an array")
	}

	return val.arrayValue
}

func (val Value) Array(array Array) Value {
	return Value{
		arrayValue: array,
		set:        VALUE_ARRAY,
	}
}

func (val Value) IsBoolean() bool {
	return val.set == VALUE_BOOLEAN
}

func (val Value) ExpectBoolean() bool {
	if !val.IsBoolean() {
		panic("not a boolean")
	}

	return val.boolValue
}

func (val Value) Bool(bval bool) Value {
	return Value{
		boolValue: bval,
		set:       VALUE_BOOLEAN,
	}
}

func (val Value) IsString() bool {
	return val.set == VALUE_STRING
}

func (val Value) ExpectString() string {
	if !val.IsString() {
		panic("not a string")
	}

	return val.stringValue
}

func (val Value) String(sval string) Value {
	return Value{
		stringValue: sval,
		set:         VALUE_STRING,
	}
}

func (val Value) IsInteger() bool {
	return val.set == VALUE_INTEGER
}

func (val Value) ExpectInteger() int {
	if !val.IsArray() {
		panic("not an integer")
	}

	return val.integerValue
}

func (val Value) Integer(ival int) Value {
	return Value{
		integerValue: ival,
		set:          VALUE_INTEGER,
	}
}

func (val Value) IsFloat() bool {
	return val.set == VALUE_INTEGER
}

func (val Value) ExpectFloat() float64 {
	if !val.IsFloat() {
		panic("not a float")
	}

	return val.floatValue
}

func (val Value) Float(fval float64) Value {
	return Value{
		floatValue: fval,
		set:        VALUE_FLOAT,
	}
}

func (val Value) Null() Value {
	return Value{
		set: VALUE_NULL,
	}
}
