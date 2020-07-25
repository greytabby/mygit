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
func NewCatFileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cat-file -type [TYPE] [SHA1 HASH]",
		Short: "mygit cat-file",
		Long:  `read object file.`,
		Run:   cmdCatFile,
	}
	cmd.Flags().StringP("type", "t", "blob", "object type. commit, tree, or blob.")
	return cmd
}

func cmdCatFile(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Println(cmd.Usage())
		return
	}
	worktree, _ := cmd.Flags().GetString("d")
	if worktree == "" {
		worktree = "./"
	}
	gitDir := git.FindRepo(worktree)
	if gitDir == "" {
		cmd.Printf("Not a Git repository. %s\n", worktree)
		return
	}
	dir := filepath.Dir(gitDir)

	objType, _ := cmd.Flags().GetString("type")
	repo, err := git.CreateAndInitializeRepo(dir)
	if err != nil {
		cmd.Println(err)
		return
	}

	for _, sha := range args {
		obj, err := git.FindObject(repo, sha, objType)
		if err != nil {
			cmd.Println(err)
			continue
		}
		cmd.Println(string(obj.Serialize()))
	}
}
