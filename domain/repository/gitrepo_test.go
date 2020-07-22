package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitDirExists(t *testing.T) {
	dir := os.TempDir()
	gitDir := filepath.Join(dir, ".git")
	err := os.Mkdir(gitDir, 0755)
	assert.NoError(t, err)
	repo, err := NewGitRepository(dir, GitRepoConfigForceMakeRepo)

	assert.NoError(t, err)
	assert.True(t, repo.GitDirExists())

	os.RemoveAll(dir)
}

func TestGitDirDoesNotExists(t *testing.T) {
	dir := os.TempDir()
	_, err := NewGitRepository(dir)

	assert.Error(t, err)

	os.RemoveAll(dir)
}

func TestGitRepoPath(t *testing.T) {
	tempDir := os.TempDir()
	repo, err := NewGitRepository(tempDir, GitRepoConfigForceMakeRepo)

	assert.NoError(t, err)
	assert.Equal(t, repo.RepositoryPath("/a/b/c.txt"), filepath.Join(tempDir, "/.git", "a", "b", "c.txt"))
	os.RemoveAll(tempDir)
}

func TestSaveRepoFile(t *testing.T) {
	cases := []struct {
		testCase string
		path     string
		content  []byte
	}{
		{testCase: "file in .git dir", path: "a.txt", content: []byte("AAAA")},
		{testCase: "file in deep dir", path: "a/b.txt", content: []byte("BBBB")},
	}

	for _, c := range cases {
		t.Run(c.testCase, func(t *testing.T) {
			tempDir := os.TempDir()
			repo, err := NewGitRepository(tempDir, GitRepoConfigForceMakeRepo)
			assert.NoError(t, err)

			err = repo.SaveRepoFile(c.path, c.content)

			assert.NoError(t, err)
			assert.FileExists(t, repo.RepositoryPath(c.path))
			os.RemoveAll(tempDir)
		})
	}
}
