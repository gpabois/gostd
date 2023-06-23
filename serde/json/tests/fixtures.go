package json_tests

import (
	"bytes"

	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/serde/json"
)

type simple struct {
	BooleanOptionalValue option.Option[bool] `serde:"boolean_optional_value"`
	IntegerValue         int                 `serde:"integer_value"`
	FloatValue           float64             `serde:"float_value"`
	StringValue          string              `serde:"string_value"`
	ArrayValue           []bool              `serde:"array_value"`
}

func fixture() simple {
	return simple{
		BooleanOptionalValue: option.Some(true),
		IntegerValue:         10,
		FloatValue:           1.10,
		StringValue:          "Hello world!",
		ArrayValue:           []bool{true, false, true},
	}
}

func encoded_fixture() []byte {
	var buf bytes.Buffer
	buf.WriteString(`{
		"boolean_optional_value": true, 
		"integer_value": 10,
		"float_value": 1.10,
		"string_value": "Hello world!",
		"array_value": [true, false, true]
	}`)
	return buf.Bytes()
}

func fixture_tokens() []json.Token {
	return []json.Token{
		json.Token{}.OpenDocument(),
		json.Token{}.String("boolean_optional_value"), json.Token{}.Colon(), json.Token{}.True(), json.Token{}.Comma(),
		json.Token{}.String("integer_value"), json.Token{}.Colon(), json.Token{}.Number("10"), json.Token{}.Comma(),
		json.Token{}.String("float_value"), json.Token{}.Colon(), json.Token{}.Number("1.10"), json.Token{}.Comma(),
		json.Token{}.String("string_value"), json.Token{}.Colon(), json.Token{}.String("Hello world!"), json.Token{}.Comma(),
		json.Token{}.String("array_value"), json.Token{}.Colon(),
		json.Token{}.OpenArray(),
		json.Token{}.True(), json.Token{}.Comma(), json.Token{}.False(), json.Token{}.Comma(), json.Token{}.True(),
		json.Token{}.CloseArray(),
		json.Token{}.CloseDocument(),
	}
}
