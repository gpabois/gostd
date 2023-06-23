package json_tests

import (
	"bytes"
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/serde/json"
)

func Test_Scanner(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("{\"val\":0}")

	scanner := json.NewScanner(&buf)

	expectedTokens := []json.Token{
		json.Token{}.OpenDocument(),
		json.Token{}.String("val"),
		json.Token{}.Colon(),
		json.Token{}.Number("0"),
		json.Token{}.CloseDocument(),
	}
	tokens := iter.CollectToSlice[[]json.Token](scanner)
}
