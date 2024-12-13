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

		fmt.Println("init called")

		if args[0] == "1" {
			lib.Init_1(args[1], args[2])
		} else if args[0] == "2" {
			lib.Init_2()
		} else if args[0] == "3" {
			lib.Init_3()
		} else {
			fmt.Println("Invalid argument")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
