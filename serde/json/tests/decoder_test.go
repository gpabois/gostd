package json_tests

import (
	"bytes"
	"testing"

	"github.com/gpabois/gostd/serde/decoder"
	"github.com/gpabois/gostd/serde/json"
	"github.com/stretchr/testify/assert"
)

func Test_Decoder(t *testing.T) {
	expectedValue := fixture()
	buf := bytes.NewBuffer(encoded_fixture())
	res := decoder.Decode[simple](json.NewDecoder(buf))

	assert.True(t, res.IsSuccess(), res.UnwrapError())
	assert.Equal(t, expectedValue, res.Expect())
}
