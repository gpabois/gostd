package json_tests

import (
	"bytes"
	"testing"

	"github.com/gpabois/gostd/serde/encoder"
	"github.com/gpabois/gostd/serde/json"
	"github.com/stretchr/testify/assert"
)

func Test_Encoder(t *testing.T) {
	v := fixture()

	var buf bytes.Buffer
	encoder.Encode(json.NewEncoder(&buf), v)

	value := buf.String()
	expectedValue := "{\"val\":0}"
	assert.Equal(t, expectedValue, value)
}
