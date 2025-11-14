package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 - calc MD5 checksum
func MD5(val []byte) string {
	hash := md5.Sum(val)
	return hex.EncodeToString(hash[:])
}
