package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	obj := NewGitBlob([]byte("test"))

	want := []byte("blob 4\x00test")
	assert.Equal(t, want, Serialize(obj))
}

func TestDeserialize(t *testing.T) {
	data := []byte("blob 4\x00test")
	want := NewGitBlob([]byte("test"))
	got, err := Deserialize(data)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
