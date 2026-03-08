package contracts

import (
	"crypto/rand"
	"encoding/hex"
)

func newEventID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
