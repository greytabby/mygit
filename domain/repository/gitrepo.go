package repository

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GitRepository struct {
	worktree string
	gitDir   string
	force    bool
}

var (
	ErrGitRepositoryAlreadyExists = errors.New(".git directory is exists. git repository is already initialized.")
	ErrNotGitRepository           = errors.New("Not a git repository.")
)

type GitRepoConfigFunc func(*GitRepository)

func NewGitRepository(path string, configFunc ...GitRepoConfigFunc) (*GitRepository, error) {
	repo := &GitRepository{
		worktree: path,
		gitDir:   filepath.Join(path, ".git"),
		force:    false,
	}

	if len(configFunc) > 0 {
		for _, f := range configFunc {
			f(repo)
		}
	}

	if !dirExists(repo.worktree) {
		errMessage := fmt.Sprintf("Worktree direcotry does not exists. %s", repo.worktree)
		return nil, errors.New(errMessage)
	}

	if !repo.force && !repo.GitDirExists() {
		return nil, ErrNotGitRepository
	}

	if repo.force && repo.GitDirExists() {
		if err := repo.SaveRepoDir(""); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

func GitRepoConfigForceMakeRepo(repo *GitRepository) {
	repo.force = true
}

func (r *GitRepository) GitDirExists() bool {
	return dirExists(r.gitDir)
}

func (r *GitRepository) RepositoryPath(path string) string {
	return filepath.Join(r.gitDir, path)
}

func (r *GitRepository) SaveRepoFile(path string, content []byte) error {
	repoPath := r.RepositoryPath(path)
	if err := os.MkdirAll(filepath.Dir(repoPath), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(repoPath, content, 0644); err != nil {
		return err
	}
	return nil
}

func (r *GitRepository) SaveRepoDir(path string) error {
	repoPath := r.RepositoryPath(path)
	return os.MkdirAll(repoPath, 0755)
}

func dirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
