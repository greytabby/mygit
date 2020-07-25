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
	"io/ioutil"
	"path/filepath"

	"github.com/greytabby/mygit/git"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
func NewHashObjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash-object -t [TYPE] [PATH]",
		Short: "write object file",
		Long:  `write object file`,
		Run:   cmdHashObject,
	}
	cmd.Flags().StringP("type", "t", "blob", "object type. commit, tree, or blob.")
	cmd.Flags().BoolP("write", "w", true, "Actually write the object into the database.")
	return cmd
}

func cmdHashObject(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
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

	objType, _ := cmd.Flags().GetString("type")
	repo, err := git.NewGitRepository(dir)
	if err != nil {
		cmd.Println(err)
		return
	}

	for _, fd := range args {
		sha, err := writeObject(fd, objType, repo)
		if err != nil {
			cmd.Println(err)
			continue
		}
		cmd.Println(fd, ":", sha)
	}
}

func writeObject(path, objType string, repo *git.GitRepository) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var fn func([]byte) git.GitObject
	switch objType {
	case "commit":
	case "tree":
	case "blob":
		fn = git.NewGitBlob
	}

	obj := fn(data)

	return git.WriteObject(repo, obj)
}
