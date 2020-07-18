package serialize

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
