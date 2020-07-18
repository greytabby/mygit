package serialize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressZlib(t *testing.T) {
	enc := CompressZlib([]byte("hello, world\n"))
	want := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
		47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}

	assert.Equal(t, want, enc)
}

func TestDeCompressZlib(t *testing.T) {
	enc := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
		47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}

	want := []byte("hello, world\n")
	got, err := DecompressZlib(enc)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}