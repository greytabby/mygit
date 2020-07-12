package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveFile(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		content    string
		permission os.FileMode
	}{
		{name: "0755", path: "a", content: "1", permission: 0755},
		{name: "0666", path: "b", content: "2", permission: 0644},
	}

	tempDir := os.TempDir()
	gitRepo := &GitRepo{
		GitDir: tempDir,
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := gitRepo.SaveFile(c.path, strings.NewReader(c.content), c.permission)
			assert.NoError(t, err)
			fp := filepath.Join(gitRepo.GitDir, c.path)
			assert.FileExists(t, fp)
			got, err := ioutil.ReadFile(fp)
			assert.NoError(t, err)
			assert.Equal(t, []byte(c.content), got)
			fi, _ := os.Stat(fp)
			assert.Equal(t, c.permission, fi.Mode())
		})
		os.RemoveAll(tempDir)
	}
}
