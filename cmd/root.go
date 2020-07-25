package cmd

import (
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
	cmd.AddCommand(NewCatFileCommand())
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
		mygit.Println(err)
		return
	}
	mygit.PersistentFlags().StringP("d", "d", dir, "git repo directory")
	if err := mygit.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
