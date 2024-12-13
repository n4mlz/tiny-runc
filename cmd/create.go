package cmd

import (
	"github.com/n4mlz/tiny-runc/lib"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <container-id>",
	Short: "create a container",
	Long: `The create command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "config.json" and a root
filesystem.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		containerID := args[0]
		bundle, err := cmd.Flags().GetString("bundle")
		if err != nil {
			panic(err)
		}

		lib.Create(containerID, bundle)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("bundle", "b", "", "path to the root of the bundle directory, defaults to the current directory")
}
