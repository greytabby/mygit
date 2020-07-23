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
	"github.com/greytabby/mygit/domain/commands"
	"github.com/spf13/cobra"
)

func NewHashObjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash-object",
		Short: "mygit hash-object",
		Long:  `write git object. only blob now`,
		Run:   cmdHashObject,
	}
	cmd.Flags().StringP("type", "t", "blob", "object type. blob, tree, or commit")
	return cmd
}

func cmdHashObject(cmd *cobra.Command, args []string) {
	req := makeHashObjectReq(cmd, args)
	commands.GitHashObject(cmd, req)
}

func makeHashObjectReq(cmd *cobra.Command, args []string) *commands.HashObjectRequest {
	gitDir, _ := cmd.Flags().GetString("git-dir")
	objType, _ := cmd.Flags().GetString("type")
	if gitDir == "" {
		gitDir = "./"
	}

	return &commands.HashObjectRequest{
		Worktree: gitDir,
		ObjType:  objType,
		Paths:    args,
	}
}
