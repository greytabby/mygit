package commands

import (
	"fmt"

	"github.com/greytabby/mygit/domain/models"
)

func GitInit(path string) {
	repo, err := models.NewGitRepository(path, models.GitRepoConfigForceMakeRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	initialDirs := []string{
		"branches",
		"objects",
		"refs/tags",
		"refs/heads",
	}
	for _, dir := range initialDirs {
		if err := repo.SaveRepoDir(dir); err != nil {
			fmt.Println(err)
			return
		}
	}

	err = repo.SaveRepoFile("HEAD", []byte("ref: refs/heads/master\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
}
