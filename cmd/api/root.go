package api

import "github.com/spf13/cobra"

func NewApiCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "api",
	}

	rootCmd.AddCommand(newApiCommand())
	return rootCmd
}
