package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "anh",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
}

func Execute(args []string) error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(NewGenerateTokenCommand())
	rootCmd.AddCommand(NewInstallationCommand())
}
