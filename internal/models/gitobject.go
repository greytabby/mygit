package models

import "io"

type GitObject interface {
	Hash() string
	PlainContent() io.Writer
	CompressedContent() io.Writer
}
