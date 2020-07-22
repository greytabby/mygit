package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/greytabby/mygit/domain/repository"
)

type GitObject interface {
	Serialize() []byte
	Deserialize([]byte)
}

func Hash(objType, data []byte) string {
	fullData := makeFullData(objType, data)
	return hash(fullData)
}

func hash(fullData []byte) string {
	sha := sha1.Sum(fullData)
	return hex.EncodeToString(sha[:])
}

func Write(repo *repository.GitRepository, objType, data []byte) (string, error) {
	fullData := makeFullData(objType, data)
	sha := hash(fullData)
	path := filepath.Join("objects", sha[:2], sha[2:])

	if err := repo.SaveRepoFile(path, CompressZlib(fullData)); err != nil {
		return "", err
	}
	return sha, nil
}

func makeFullData(objType, data []byte) []byte {
	size := strconv.Itoa(len(data))
	return bytes.Join([][]byte{objType, []byte(" "), []byte(size), []byte("\x00"), data}, []byte(""))
}

func CompressZlib(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

func DecompressZlib(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
