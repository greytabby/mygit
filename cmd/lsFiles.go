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
func NewLsFilesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-files",
		Short: "print list of files in index",
		Long:  `print list of files in index`,
		Run:   cmdLsFiles,
	}
	return cmd
}

func cmdLsFiles(cmd *cobra.Command, args []string) {
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
	if err != nil {
		cmd.Println(err)
		return
	}

	for _, entry := range index.Entries {
		stage := (entry.Flags >> 12) & 3
		cmd.Printf("%s %s %d\t%s\n", entry.Mode.Perm().String(), entry.ObjectID, int(stage), entry.FilePath)
	}
}
