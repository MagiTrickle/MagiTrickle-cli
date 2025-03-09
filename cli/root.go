package cli

import (
	"github.com/spf13/cobra"
)

// rootCmd is the main command (entry point)
var rootCmd = &cobra.Command{
	Use:   "magitrickle",
	Short: "A CLI tool for managing MagiTrickle via UNIX socket",
	Long: `MagiTrickle CLI communicates with the MagiTrickle API through a 
UNIX socket (http://unix). It supports commands for system hooks, interfaces, 
configuration, group management, and more.
	
Try these commands:
  magitrickle system interfaces
  magitrickle group list
  magitrickle group create --name=MyGroup --interface=br0
`,
}

// Execute launches the root command
func Execute() {
	_ = rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(systemCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(ruleCmd)
}
