package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/greytabby/mygit/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestHashObject(t *testing.T) {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	err := ioutil.WriteFile(tempFile, []byte("test"), 0755)
	GitInit(tempDir)

	_, err = models.NewGitRepository(tempDir, models.GitRepoConfigForceMakeRepo)
	assert.NoError(t, err)
	GitHashObject(tempDir, true, "test.txt")

	os.RemoveAll(tempDir)
}
