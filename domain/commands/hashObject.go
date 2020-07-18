package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/greytabby/mygit/domain/models"
	"github.com/greytabby/mygit/domain/serialize"
)

func GitHashObject(gitDir string, w bool, path string) {
	repo, err := models.NewGitRepository(gitDir)
	if err != nil {
		fmt.Println(err)
	}

	file := filepath.Join(gitDir, path)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("No such file or directory. %s\n", path)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	obj := models.NewGitBlob(data)
	serializedData := models.Serialize(obj)
	hash := serialize.Hash(serializedData)
	compressedData := serialize.CompressZlib(serializedData)
	objPath := filepath.Join("objects", hash[:2], hash[2:])
	repo.SaveRepoFile(objPath, compressedData)
}
