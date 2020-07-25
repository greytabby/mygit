/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"path/filepath"

	"github.com/greytabby/mygit/git"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
func NewLsTreeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-tree [TREE OBJECT]",
		Short: "listing tree object entries",
		Long:  `listing tree object entries`,
		Run:   cmdLsTree,
	}
	return cmd
}

func cmdLsTree(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Println(cmd.Usage())
		return
	}
	worktree, _ := cmd.Flags().GetString("d")
	if worktree == "" {
		worktree = "./"
	}
	gitDir, err := git.FindRepo(worktree)
	if err != nil {
		cmd.Println(err)
		return
	}
	dir := filepath.Dir(gitDir)

	repo, err := git.NewGitRepository(dir)
	if err != nil {
		cmd.Println(err)
		return
	}

	obj, err := git.ReadObject(repo, args[0])
	if err != nil {
		cmd.Println(err)
		return
	}

	for _, entry := range obj.(*git.GitTree).Entries {
		mode := entry.Mode.Perm().String()
		obj, err := git.ReadObject(repo, entry.Sha)
		if err != nil {
			cmd.Println(err)
			return
		}
		cmd.Printf("%s %s %s\t%s\n", mode, string(obj.Type()), entry.Sha, entry.Path)
	}
}
