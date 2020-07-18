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
	"github.com/spf13/cobra"

	"github.com/greytabby/mygit/domain/commands"
)

func NewHashObjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash-object",
		Short: "mygit hash-object",
		Long:  `write git object. only blob now`,
		Run:   cmdHashObject,
	}
	cmd.Flags().BoolP("w", "w", true, "write the object into the object database")
	cmd.Flags().StringP("path", "p", nil, "process file as it were from this path")
	return cmd
}

func cmdHashObject(cmd *cobra.Command, args []string) {
	gitDir, _ := cmd.Flags().GetString("git-dir")
	w := cmd.Flags().GetBool("w")
	path := cmd.Flags().GetString("path")
	if gitDir == "" {
		gitDir = "./"
	}
	commands.GitHashObject(gitDir, w, path)
}
