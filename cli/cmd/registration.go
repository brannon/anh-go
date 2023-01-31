package cmd

import "github.com/spf13/cobra"

func NewRegistrationCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use: "registration",
	}

	cobraCmd.AddCommand(
		NewRegistrationGetCommand(),
		NewRegistrationListCommand(),
	)

	return cobraCmd
}
