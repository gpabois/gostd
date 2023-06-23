package json_tests

import (
	"bytes"
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/serde/json"
	"github.com/stretchr/testify/assert"
)

func Test_Scanner(t *testing.T) {
	buf := bytes.NewBuffer(encoded_fixture())
	scanner := json.NewScanner(buf)

	expectedTokens := fixture_tokens()
	tokens := iter.CollectToSlice[[]json.Token, json.Token](scanner)

	assert.Equal(t, expectedTokens, tokens)
}
