package cmd

import "github.com/spf13/cobra"

func NewInstallationCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use: "installation",
	}

	cobraCmd.AddCommand(
		NewInstallationGetCommand(),
		NewInstallationListCommand(),
	)

	return cobraCmd
}
