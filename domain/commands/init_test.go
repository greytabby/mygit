package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitInit(t *testing.T) {
	tempDir := os.TempDir()
	GitInit(tempDir)

	initialDirs := []string{
		"branches",
		"objects",
		"refs/tags",
		"refs/heads",
	}
	for _, dir := range initialDirs {
		path := filepath.Join(tempDir, ".git", dir)
		assert.DirExists(t, path)
	}
	path := filepath.Join(tempDir, ".git", "HEAD")
	assert.FileExists(t, path)

	os.RemoveAll(tempDir)
}