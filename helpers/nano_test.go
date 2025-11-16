package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDWithLength(t *testing.T) {
	// when
	result := ID(11)

	// then
	assert.NotEmpty(t, result)
	assert.Len(t, result, 11)
}

func TestIDWithoutLength(t *testing.T) {
	// when
	result := ID()

	// then
	assert.NotEmpty(t, result)
	assert.Len(t, result, 11)
}

func TestIDInvalidLength(t *testing.T) {
	// when
	result := ID(-5)

	// then
	assert.NotEmpty(t, result)
	assert.Len(t, result, 11)
}
