package rand

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gpabois/gostd/result"
)

func RandomHex(n int) result.Result[string] {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return result.Failed[string](err)
	}
	return result.Success(hex.EncodeToString(bytes))
}
