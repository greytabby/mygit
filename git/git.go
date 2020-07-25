package git

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type IndexEntry struct {
	ctime    int64
	mtime    int64
	dev      int32
	inot     int32
	mode     os.FileMode
	uid      int
	gid      int
	fileSize int32
	objectID []byte
	flags    int16
	filePath string
}

type GitRepository struct {
	Worktree string
	GitDir   string
}

type GitObject interface {
	Serialize() []byte
	Deserialize([]byte)
	Type() []byte
}

// NewGitRepository return `GitRepository`
// return error if ".git" dir is not exists in path.
func NewGitRepository(path string) (*GitRepository, error) {
	return newRepo(path, false)
}

// NewGitRepository return `GitRepository`
// force create if path is not Git directory.
func NewGitRepositoryForceCreate(path string) (*GitRepository, error) {
	return newRepo(path, true)
}

func newRepo(path string, force bool) (*GitRepository, error) {
	repo := &GitRepository{
		Worktree: path,
		GitDir:   filepath.Join(path, ".git"),
	}
	info, err := os.Stat(repo.GitDir)
	if err != nil && os.IsExist(err) {
		return nil, err
	}
	if !(force || info.IsDir()) {
		return nil, fmt.Errorf("Not a Git repository. %s", path)
	}
	return repo, nil
}

// RepoPath return path joined ".git" directory
func (gr *GitRepository) RepoPath(path string) string {
	return filepath.Join(gr.GitDir, path)
}

// SaveRepoFile save file to path which joined ".git".
func (gr *GitRepository) SaveRepoFile(path string, data []byte) error {
	savePath := filepath.Join(gr.GitDir, path)
	baseDir := filepath.Dir(savePath)
	if err := os.MkdirAll(baseDir, repoDirPerm()); err != nil {
		return err
	}
	if err := ioutil.WriteFile(savePath, data, repoFilePerm()); err != nil {
		return err
	}
	return nil
}

func (gr *GitRepository) CreateRepoDir(path string) error {
	repoDir := filepath.Join(gr.GitDir, path)
	if err := os.MkdirAll(repoDir, repoDirPerm()); err != nil {
		return err
	}
	return nil
}

// Create a new repository at path
func CreateAndInitializeRepo(path string) (*GitRepository, error) {
	repo, err := NewGitRepositoryForceCreate(path)
	if err != nil {
		return nil, err
	}
	dirs := []string{
		"objects",
		"branches",
		"refs/tags",
		"refs/heads",
	}
	for _, dir := range dirs {
		if err := repo.CreateRepoDir(dir); err != nil {
			return nil, err
		}
	}

	// .git/description
	repo.SaveRepoFile("description", []byte("Unnamed repository; edit this file 'description' to name the repository.\n"))

	// .git/HEAD
	repo.SaveRepoFile("HEAD", []byte("ref: refs/heads/master\n"))

	return repo, nil
}

// FindRepo return git repository path
// if .git directory is not exist in path,
// find parent directory recursively.
func FindRepo(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	repoDir := filepath.Join(absPath, ".git")
	if _, err := os.Stat(repoDir); !os.IsNotExist(err) {
		return repoDir, nil
	}

	parent := filepath.Join(absPath, "../")
	if parent == absPath {
		return "", fmt.Errorf("Git repository not found. %s\n", path)
	}

	return FindRepo(parent)
}

func repoFilePerm() os.FileMode {
	return os.FileMode(0644)
}

func repoDirPerm() os.FileMode {
	return os.FileMode(0755)
}

type (
	GitCommit struct{}
	GitTree   struct{}
	GitBlob   struct {
		Data []byte
	}
)

func NewGitBlob(data []byte) GitObject {
	o := new(GitBlob)
	o.Deserialize(data)
	return o
}

func (o *GitBlob) Serialize() []byte {
	return o.Data
}

func (o *GitBlob) Deserialize(data []byte) {
	o.Data = data
}

func (o *GitBlob) Type() []byte {
	return []byte("blob")
}

// ReadObject retrun a GitObject whose exact type depends on the object.
func ReadObject(repo *GitRepository, sha string) (GitObject, error) {
	path := filepath.Join("objects", sha[:2], sha[2:])
	encData, err := ioutil.ReadFile(repo.RepoPath(path))
	if err != nil {
		return nil, err
	}
	data, err := decompressZlib(encData)
	if err != nil {
		return nil, err
	}

	// Read Object Type
	x := bytes.IndexByte(data, ' ')
	objType := string(data[:x])

	// Read Object size
	y := bytes.IndexByte(data, '\x00')
	size, err := byte2Int(data[x+1 : y])
	if size != len(data)-y-1 {
		return nil, fmt.Errorf("Malformed object %d: bad length", size)
	}

	var fn func([]byte) GitObject
	switch objType {
	case "commit":
	case "tree":
	case "blob":
		fn = NewGitBlob
	}

	return fn(data[y+1:]), nil
}

func WriteObject(repo *GitRepository, obj GitObject) (string, error) {
	sha, data := HashObject(obj)
	if err := repo.SaveRepoFile(filepath.Join("objects", sha[:2], sha[2:]), compressZlib(data)); err != nil {
		return "", err
	}
	return sha, nil
}

func FindObject(repo *GitRepository, sha, objType string) (GitObject, error) {
	if len(sha) < 3 {
		return nil, fmt.Errorf("hash prefix must be 2 or more charcters.\n")
	}
	var files []string
	objDir := repo.RepoPath(filepath.Join("objects", sha[:2]))
	err := filepath.Walk(objDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), sha[2:]) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		hash := sha[:2] + filepath.Base(file)
		obj, err := ReadObject(repo, hash)
		if err != nil {
			return nil, err
		}
		if string(obj.Type()) == objType {
			return obj, nil
		}
	}
	return nil, fmt.Errorf("Object not found. %s\n", sha)
}

// HashObject return object hash and serialized data
func HashObject(obj GitObject) (string, []byte) {
	data := obj.Serialize()
	size := []byte(strconv.Itoa(len(data)))
	result := bytes.Join([][]byte{obj.Type(), []byte(" "), size, []byte("\x00"), data}, []byte(""))
	sha := hash(result)
	return sha, result
}

func hash(data []byte) string {
	sha := sha1.Sum(data)
	return hex.EncodeToString(sha[:])
}

func compressZlib(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

func decompressZlib(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func byte2Int(b []byte) (int, error) {
	return strconv.Atoi(string(b))
}
