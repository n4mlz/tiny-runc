package cmd

import (
	"fmt"
	"os"

	"github.com/n4mlz/tiny-runc/lib"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:    "init",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Invalid number of arguments")
			os.Exit(1)
		}

		containerID := args[1]

		if args[0] == "1" {
			pipeFromParent := args[2]
			pipeToParent := args[3]
			lib.Init_1(containerID, pipeFromParent, pipeToParent)
		} else if args[0] == "2" {
			lib.Init_2(containerID)
		} else if args[0] == "3" {
			lib.Init_3(containerID)
		} else {
			fmt.Println("Invalid argument")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
