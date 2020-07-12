package models

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GitRepo struct {
	GitDir string
}

func (d *GitRepo) SaveFile(path string, content io.Reader, permission os.FileMode) error {
	savePath := filepath.Join(d.GitDir, path)
	data, err := ioutil.ReadAll(content)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(savePath, data, permission); err != nil {
		return err
	}
	return nil
}
