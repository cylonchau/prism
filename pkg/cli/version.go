package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version info (injected at build time)
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, git commit, and build date of prism.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Prism Version: %s\n", Version)
		fmt.Printf("Git Commit:    %s\n", GitCommit)
		fmt.Printf("Build Date:    %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
