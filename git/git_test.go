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

	found, err := FindRepo(dirA)
	assert.NoError(t, err)
	assert.Equal(t, gitDir, found)
	notFound, err := FindRepo(dirB)
	assert.Error(t, err)
	assert.Equal(t, "", notFound)

	os.RemoveAll(temp)
}

func TestHashObject(t *testing.T) {
	obj := NewGitBlob([]byte("test\n"))
	sha, data := HashObject(obj)
	wantData := []byte("blob 5\x00test\n")
	assert.Equal(t, "9daeafb9864cf43055ae93beb0afd6c7d144bfa4", sha)
	assert.Equal(t, wantData, data)
}

func TestWriteObject(t *testing.T) {
	temp := os.TempDir()
	repo, _ := CreateAndInitializeRepo(temp)
	obj := NewGitBlob([]byte("test\n"))

	sha, err := WriteObject(repo, obj)

	assert.NoError(t, err)
	assert.Equal(t, "9daeafb9864cf43055ae93beb0afd6c7d144bfa4", sha)
	objPath := repo.RepoPath(filepath.Join("objects", sha[:2], sha[2:]))
	assert.FileExists(t, objPath)

	os.RemoveAll(temp)
}

func TestReadObject(t *testing.T) {
	temp := os.TempDir()
	repo, _ := CreateAndInitializeRepo(temp)
	obj := NewGitBlob([]byte("test\n"))
	sha, err := WriteObject(repo, obj)
	assert.NoError(t, err)

	got, err := ReadObject(repo, sha)
	assert.NoError(t, err)
	assert.Equal(t, obj.Serialize(), got.Serialize())

	os.RemoveAll(temp)
}

func TestFindObject(t *testing.T) {
	temp := os.TempDir()
	repo, _ := CreateAndInitializeRepo(temp)
	obj := NewGitBlob([]byte("test\n"))
	sha, _ := WriteObject(repo, obj)

	got, err := FindObject(repo, sha[:3], "blob")
	assert.NoError(t, err)
	assert.Equal(t, obj.Serialize(), got.Serialize())

	os.RemoveAll(temp)
}
