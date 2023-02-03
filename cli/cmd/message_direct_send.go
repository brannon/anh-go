package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/brannon/anh-go"
	"github.com/spf13/cobra"
)

func NewMessageDirectSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "direct-send",
		RunE: func(c *cobra.Command, args []string) error {
			deviceToken, _ := c.PersistentFlags().GetString("device-token")
			platform, _ := c.PersistentFlags().GetString("platform")
			stdin, _ := c.PersistentFlags().GetBool("stdin")
			bodyString, _ := c.PersistentFlags().GetString("body")

			var err error
			var notification anh.Notification

			switch platform {
			case "apple", "apns":
				appleNotification := &anh.AppleNotification{}

				if stdin {
					err = appleNotification.SetBodyFromReader(os.Stdin)
				} else if bodyString != "" {
					err = appleNotification.SetBodyFromString(bodyString)
				} else {
					err = fmt.Errorf("body is required")
				}

				if err != nil {
					return err
				}

				notification = appleNotification

			default:
				return fmt.Errorf("invalid platform: %s", platform)
			}

			if deviceToken == "" {
				return fmt.Errorf("device token is required")
			}

			hubName, connectionString, err := getHubNameAndConnectionString()
			if err != nil {
				return err
			}

			client, err := anh.NewClient(hubName, anh.WithConnectionString(connectionString))
			if err != nil {
				return err
			}

			client.VerboseLogger = getVerboseLogger()

			notificationResult, err := client.SendDirectNotification(context.Background(), notification, deviceToken)
			if err != nil {
				return err
			}

			fmt.Printf("Notification ID: %s\n", notificationResult.GetNotificationId())

			return nil
		},
	}

	cmd.PersistentFlags().Bool("stdin", false, "Read message from stdin")
	cmd.PersistentFlags().String("platform", "", "Platform (apple, android)")
	cmd.PersistentFlags().String("device-token", "", "Device token")
	cmd.PersistentFlags().String("body", "", "Message body (JSON)")

	return cmd
}
