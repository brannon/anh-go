package cmd

import (
	"context"
	"fmt"

	"github.com/brannon/anh-go/anh"
	"github.com/spf13/cobra"
)

func NewInstallationGetCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:  "get",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			installationId := args[0]

			hubName, connectionString, err := getHubNameAndConnectionString()
			if err != nil {
				return err
			}

			client, err := anh.NewClient(hubName, connectionString)
			if err != nil {
				return err
			}

			client.Logger = getLogger()

			installation, err := client.GetInstallation(context.Background(), installationId)
			if err != nil {
				return err
			}

			fmt.Println(installation.GetRawData().FormattedString())

			return nil
		},
	}

	return cobraCmd
}
