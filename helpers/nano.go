package helpers

import (
	"fmt"
	"math/rand"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

const nanoIDLen = 11

func ID(length ...int) string {
	genLength := nanoIDLen
	if len(length) > 0 {
		genLength = length[0]
	}

	id, err := gonanoid.ID(genLength)
	if err != nil {
		fmt.Println("gen nano ID:", err)
		if genLength < 0 {
			genLength = nanoIDLen
		}

		id = GetRandomString(genLength)
	}
	return id
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetRandomString - superfast string generation.
// 139 ns/op   32 B/op   2 allocs/op.
// topic: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func GetRandomString(length int) string {
	var rndSource = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, rndSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rndSource.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
