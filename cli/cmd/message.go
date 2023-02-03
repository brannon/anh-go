package cmd

import "github.com/spf13/cobra"

func NewMessageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "message",
	}

	cmd.AddCommand(NewMessageDirectSendCommand())

	return cmd
}
