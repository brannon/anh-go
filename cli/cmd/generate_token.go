package cmd

import (
	"fmt"
	"time"

	"github.com/brannon/anh-go"
	"github.com/spf13/cobra"
)

func NewGenerateTokenCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:  "generate-token",
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			var connectionString string

			if len(args) > 0 {
				connectionString = args[0]
			} else {
				var err error
				_, connectionString, err = getHubNameAndConnectionString()
				if err != nil {
					return err
				}
			}

			cs, err := anh.ParseConnectionString(connectionString)
			if err != nil {
				return err
			}

			tokenProvider := anh.NewSasTokenProvider(cs.KeyName, cs.Key)

			token, _, err := tokenProvider.GenerateSasToken(cs.Endpoint, time.Now().Add(1*time.Hour))
			if err != nil {
				return err
			}

			fmt.Println(token)
			return nil
		},
	}

	return cobraCmd
}
