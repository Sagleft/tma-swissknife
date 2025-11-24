package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxHash64Base32(t *testing.T) {
	// given
	data := []byte("test")

	// when
	result := XxHash64Base32(data)

	// then
	assert.NotEmpty(t, result)
}
