package worker

import "github.com/spf13/cobra"

func NewWorkerCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "worker",
	}

	rootCmd.AddCommand(newPendingMessagesWorkerCommands())
	rootCmd.AddCommand(newSendMessageWorkerCommands())
	return rootCmd
}
