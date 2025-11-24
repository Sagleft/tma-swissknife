package helpers

import (
	"crypto/md5"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/cespare/xxhash/v2"
)

// MD5 - calc MD5 checksum
func MD5(val []byte) string {
	hash := md5.Sum(val)
	return hex.EncodeToString(hash[:])
}

func XxHash64Base32(data []byte) string {
	hash := xxhash.Sum64(data)
	// Преобразуем uint64 в байты в порядке little-endian
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, hash)
	// Кодируем байты в base32 без padding
	return strings.ToLower(
		base32.StdEncoding.WithPadding(base32.NoPadding).
			EncodeToString(buf),
	)
}
