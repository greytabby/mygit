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
	"os"
	"path/filepath"

	"github.com/greytabby/mygit/git"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [FILE...]",
		Short: "staging file",
		Long:  `staging file`,
		Run:   cmdAdd,
	}
	return cmd
}

func cmdAdd(cmd *cobra.Command, args []string) {
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

	index, err := git.ReadIndex(repo)
	if err != nil && os.IsNotExist(err) {
		index = &git.GitIndex{}
		err = nil
	}
	if err != nil {
		cmd.Println(err)
		return
	}

	var entries []*git.IndexEntry
	for _, e := range index.Entries {
		exist := false
		for _, path := range args {
			if e.FilePath == path {
				exist = true
				break
			}
		}
		if !exist {
			entries = append(entries, e)
		}
	}

	for _, path := range args {
		sha, err := writeObject(path, "blob", repo)
		if err != nil {
			cmd.Println(err)
			return
		}
		info, err := os.Stat(path)
		if err != nil {
			cmd.Println(err)
			return
		}
		e := git.NewIndexEntry(info, path, sha)
		entries = append(entries, e)
	}

	// write index file
	newIndex := &git.GitIndex{Entries: entries}
	err = git.WriteIndex(repo, newIndex)
	if err != nil {
		cmd.Println(err)
		return
	}
}
