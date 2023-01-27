package cmd

import (
	"context"
	"fmt"

	"github.com/brannon/anh-go/anh"
	"github.com/spf13/cobra"
)

func NewInstallationListCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			hubName, connectionString, err := getHubNameAndConnectionString()
			if err != nil {
				return err
			}

			client, err := anh.NewClient(hubName, connectionString)
			if err != nil {
				return err
			}

			client.Logger = getLogger()

			collection, err := client.ListInstallations(context.Background())
			if err != nil {
				return err
			}

			for collection.HasItems() {
				for _, installation := range collection.Items() {
					fmt.Println(installation.GetRawData().String())
				}

				err = collection.NextPage(context.Background())
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cobraCmd
}
