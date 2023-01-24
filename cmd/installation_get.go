package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/brannon/anh-go/anh"
	"github.com/spf13/cobra"
)

func NewInstallationGetCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:  "get",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			installationId := args[0]

			client, err := anh.NewClient(HubName, ConnectionString)
			if err != nil {
				return err
			}

			log.SetOutput(os.Stderr)
			client.Logger = log.Default()

			installation, err := client.GetInstallation(context.Background(), installationId)
			if err != nil {
				return err
			}

			fmt.Println(installation.GetRawData().PrettyString())

			return nil
		},
	}

	return cobraCmd
}
