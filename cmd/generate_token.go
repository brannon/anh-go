package cmd

import (
	"fmt"
	"time"

	"github.com/brannon/anh-go/anh"
	"github.com/spf13/cobra"
)

func NewGenerateTokenCommand() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:  "generate-token",
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			connectionString := args[0]

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
