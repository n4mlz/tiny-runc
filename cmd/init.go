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
		bundle, err := cmd.Flags().GetString("bundle")
		if err != nil {
			panic(err)
		}

		if args[0] == "1" {
			pipeFromParent := args[2]
			pipeToParent := args[3]
			lib.Init_1(containerID, bundle, pipeFromParent, pipeToParent)
		} else if args[0] == "2" {
			lib.Init_2(containerID, bundle)
		} else if args[0] == "3" {
			lib.Init_3(containerID, bundle)
		} else {
			fmt.Println("Invalid argument")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("bundle", "b", "", "path to the root of the bundle directory, defaults to the current directory")
}
