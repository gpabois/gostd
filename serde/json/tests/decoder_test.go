package json_tests

import (
	"bytes"
	"testing"

	"github.com/gpabois/gostd/serde/decoder"
	"github.com/gpabois/gostd/serde/json"
	"github.com/stretchr/testify/assert"
)

func Test_Decoder(t *testing.T) {
	v := simpleStruct{}

	var buf bytes.Buffer
	buf.WriteString("{\"val\":0}")
	res := decoder.Decode[simpleStruct](json.NewDecoder(&buf))

	assert.True(t, res.IsSuccess(), res.UnwrapError())
	assert.Equal(t, v, res.Expect())
}
