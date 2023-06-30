package norm_tests

import (
	"github.com/gpabois/gostd/option"
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

func encoded_fixture() map[string]any {
	return map[string]any{
		"boolean_optional_value": "true",
		"integer_value":          "10",
		"float_value":            "1.10",
		"string_value":           "Hello world!",
		"array_value":            []string{"true", "false", "true"},
	}
}
