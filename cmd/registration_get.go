package cmd

import (
	"context"
	"fmt"

	"github.com/brannon/anh-go/anh"
	"github.com/spf13/cobra"
)

func NewRegistrationGetCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:  "get",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			registrationId := args[0]

			hubName, connectionString, err := getHubNameAndConnectionString()
			if err != nil {
				return err
			}

			client, err := anh.NewClient(hubName, connectionString)
			if err != nil {
				return err
			}

			client.Logger = getLogger()

			registration, err := client.GetRegistration(context.Background(), registrationId)
			if err != nil {
				return err
			}

			fmt.Printf("ID: %s; Platform: %s, Tags: %v\n", registration.GetRegistrationId(), registration.GetPlatform(), registration.GetTags())

			return nil
		},
	}

	return cobraCmd
}
