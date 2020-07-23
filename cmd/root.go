package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewMygitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mygit",
		Short: "mygit is reinventing the wheel",
		Long:  "mygit is git implementetion in golang.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(NewInitCommand())
	cmd.AddCommand(NewHashObjectCommand())
	return cmd
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mygit v0.0.1")
	},
}

func Execute() {
	mygit := NewMygitCommand()
	dir, err := os.Getwd()
	if err != nil {
		errors.New("Cannot get current directory name")
	}
	mygit.PersistentFlags().StringP("git-dir", "d", dir, "git repo directory")
	if err := mygit.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
