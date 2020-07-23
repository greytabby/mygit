package commands

import (
	"github.com/greytabby/mygit/domain/objects"
	"github.com/greytabby/mygit/domain/repository"
	"github.com/spf13/cobra"
)

type HashObjectRequest struct {
	Worktree string
	ObjType  string
	Paths    []string
}

func GitHashObject(cmd *cobra.Command, req *HashObjectRequest) {
	repo, err := repository.NewGitRepository(req.Worktree)
	if err != nil {
		cmd.Println(err)
	}

	for _, path := range req.Paths {
		switch req.ObjType {
		case "blob":
			hashObjectBlob(cmd, repo, path)
		case "tree":
			cmd.Println("Not Imprement!!")
		case "Commit":
			cmd.Println("Not Imprement!!")
		default:
			cmd.PrintErrf("Unknown object type. %s\n", req.ObjType)
		}
	}
}

func hashObjectBlob(cmd *cobra.Command, repo *repository.GitRepository, path string) {
	data, err := repo.ReadWorkTreeFile(path)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	blob := objects.NewBlob(data)

	hash, err := objects.Write(repo, blob.Type(), blob.Serialize())
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	cmd.Println(hash)
}
