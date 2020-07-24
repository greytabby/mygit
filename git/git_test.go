package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveRepoFile(t *testing.T) {
	temp := os.TempDir()
	repo, err := CreateAndInitializeRepo(temp)
	assert.NoError(t, err)

	path := "dir/file"
	data := []byte("test data \x00\n")
	err = repo.SaveRepoFile(path, data)

	assert.FileExists(t, repo.RepoPath(path))
	wantPerm := os.FileMode(0644)
	info, err := os.Stat(repo.RepoPath(path))
	assert.NoError(t, err)
	assert.Equal(t, wantPerm, info.Mode().Perm())
	gotData, err := ioutil.ReadFile(repo.RepoPath(path))
	assert.NoError(t, err)
	assert.Equal(t, data, gotData)

	os.RemoveAll(temp)
}

func TestCreateRepoDir(t *testing.T) {
	temp := os.TempDir()
	repo, err := CreateAndInitializeRepo(temp)
	assert.NoError(t, err)

	path := "dir"
	err = repo.CreateRepoDir(path)

	assert.DirExists(t, repo.RepoPath(path))
	wantPerm := os.FileMode(0755)
	info, err := os.Stat(repo.RepoPath(path))
	assert.NoError(t, err)
	assert.Equal(t, wantPerm, info.Mode().Perm())

	os.RemoveAll(temp)
}

func TestFindRepo(t *testing.T) {
	temp := os.TempDir()
	dirA := filepath.Join(temp, "test/aaa/bbb")
	dirB := filepath.Join(temp, "Nogit")
	_ = os.MkdirAll(dirA, 0755)
	_ = os.MkdirAll(dirB, 0755)
	worktree := filepath.Join(temp, "test")
	gitDir := filepath.Join(worktree, ".git")
	_, err := CreateAndInitializeRepo(worktree)
	assert.NoError(t, err)

	assert.Equal(t, gitDir, FindRepo(dirA))
	assert.Equal(t, "", FindRepo(dirB))

	os.RemoveAll(temp)
}