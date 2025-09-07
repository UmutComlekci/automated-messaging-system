package cmd

import (
	"github.com/spf13/cobra"
	"github.com/umutcomlekci/automated-messaging-system/cmd/api"
	"github.com/umutcomlekci/automated-messaging-system/cmd/worker"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "",
	}

	rootCmd.AddCommand(api.NewApiCommands())
	rootCmd.AddCommand(worker.NewWorkerCommands())
	return rootCmd
}
