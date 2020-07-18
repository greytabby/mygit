package models

import "strconv"

type GitBlob struct {
	data []byte
	size int
}

func NewGitBlob(data []byte) *GitBlob {
	return &GitBlob{
		data: data,
		size: len(data),
	}
}

func (o *GitBlob) Type() []byte {
	return []byte("blob")
}

func (o *GitBlob) Data() []byte {
	return o.data
}

func (o *GitBlob) Size() []byte {
	size := strconv.Itoa(o.size)
	return []byte(size)
}
