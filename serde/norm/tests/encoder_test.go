package norm_tests

import (
	"testing"

	"github.com/gpabois/gostd/serde/encoder"
	"github.com/gpabois/gostd/serde/norm"
	"github.com/stretchr/testify/assert"
)

func Test_Encoder(t *testing.T) {
	expectedValue := encoded_fixture()

	e := norm.NewEncoder()
	encoder.Encode(e, fixture())
	assert.Equal(t, expectedValue, e.Encoded())
}
