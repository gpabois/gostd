package norm_tests

import (
	"testing"

	"github.com/gpabois/gostd/serde/decoder"
	"github.com/gpabois/gostd/serde/norm"
	"github.com/stretchr/testify/assert"
)

func Test_Decoder(t *testing.T) {
	expectedValue := fixture()
	normValue := encoded_str_fixture()

	d := norm.NewDecoder(normValue)
	res := decoder.Decode[simple](d)

	assert.True(t, res.IsSuccess(), res.UnwrapError())
	assert.Equal(t, expectedValue, res.Expect())
}
