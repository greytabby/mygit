package models

import (
	"bytes"
	"errors"
	"fmt"
)

type GitObject interface {
	Type() []byte
	Data() []byte
	Size() []byte
}

func Serialize(o GitObject) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(o.Type())
	buf.Write([]byte(" "))
	buf.Write(o.Size())
	buf.Write([]byte{0})
	buf.Write(o.Data())
	return buf.Bytes()
}

func Deserialize(data []byte) (GitObject, error) {
	x := bytes.IndexByte(data, ' ')
	objType := data[:x]
	y := bytes.Index(data, []byte{0})
	_ = data[x:y]
	objData := data[y+1:]

	switch string(objType) {
	case "blob":
		return NewGitBlob(objData), nil
	default:
		errMessage := fmt.Sprintf("Unknown object type. %s", objType)
		return nil, errors.New(errMessage)
	}
}
