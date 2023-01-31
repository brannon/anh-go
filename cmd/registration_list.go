package cmd

import (
	"context"
	"fmt"

	"github.com/brannon/anh-go/anh"
	"github.com/spf13/cobra"
)

func NewRegistrationListCommand() *cobra.Command {
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

			collection, err := client.ListRegistrations(context.Background())
			if err != nil {
				return err
			}

			for collection.HasItems() {
				for _, registration := range collection.Items() {
					fmt.Printf("ID: %s; Platform: %s, Tags: %v\n", registration.GetRegistrationId(), registration.GetPlatform(), registration.GetTags())
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
