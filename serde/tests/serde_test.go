package serde_tests

import (
	"testing"

	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/serde"
	"github.com/stretchr/testify/assert"
)

type subTestStruct struct {
	El0 int
	El1 []bool
	El2 string
	El3 map[string]string
}
type testStruct struct {
	OptValue    option.Option[string] `serde:"opt_value"`
	StructValue subTestStruct         `serde:"struct_value"`
}

func fixture() testStruct {
	return testStruct{
		OptValue: option.Some("test"),
		StructValue: subTestStruct{
			El0: 10,
			El1: []bool{true, false, true},
		},
	}
}

func zero_fixture() testStruct {
	return testStruct{}
}

func testSerde(t *testing.T, contentType string) {
	expectedVal := fixture()

	serRes := serde.Serialize(expectedVal, contentType)
	assert.True(t, serRes.IsSuccess(), serRes.UnwrapError())

	deserRes := serde.Deserialize[testStruct](serRes.Expect(), contentType)
	assert.True(t, deserRes.IsSuccess(), deserRes.UnwrapError())

	val := deserRes.Expect()
	assert.Equal(t, expectedVal, val)
}

func testZeroValue(t *testing.T, contentType string) {
	expectedVal := zero_fixture()

	serRes := serde.Serialize(expectedVal, contentType)
	assert.True(t, serRes.IsSuccess(), serRes.UnwrapError())

	deserRes := serde.Deserialize[testStruct](serRes.Expect(), contentType)
	assert.True(t, deserRes.IsSuccess(), deserRes.UnwrapError())

	val := deserRes.Expect()
	assert.Equal(t, expectedVal, val)
}

func Test_Serde_Json(t *testing.T) {
	testSerde(t, "application/json")
	testZeroValue(t, "application/json")
}
