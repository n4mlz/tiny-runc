package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tiny-runc",
	Short: "Open Container Initiative runtime",
	Long: `tiny-runc (WIP)

This is a lightweight version of the low-level container runtime, runc (WIP).

It aims to serve as a reference material for understanding how container runtimes operate, rather than focusing on practical use.

The following goals are set:
- Compliance with the essential parts of the OCI Runtime Specification
- **Rootless operation**
- Simple and easy-to-understand code
- Serve as a reference for creating custom container runtimes

Conversely, the following are not within the scope of the goals:
- Full compliance with the OCI Runtime Specification
- Advanced security considerations
- Use in production environments
- Integration with high-level container runtimes

License:
This project is licensed under the MIT License.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
