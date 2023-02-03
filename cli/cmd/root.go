package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "anh",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute(args []string) error {
	log.SetOutput(os.Stderr)

	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Configure config file search paths
	viper.SetConfigName(".anh")
	viper.SetConfigType("yaml")
	dir, err := homedir.Dir()
	if err == nil {
		viper.AddConfigPath(dir)
		viper.AddConfigPath(path.Join(dir, ".anh"))
	} else {
		log.Printf("WARNING: could not determine user's home directory: %v", err)
	}
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	// Configure ENV var config mapping
	viper.SetEnvPrefix("anh")
	viper.BindEnv("connection_string")
	viper.BindEnv("hub_name")

	// Run the command
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(NewGenerateTokenCommand())
	rootCmd.AddCommand(NewInstallationCommand())
	rootCmd.AddCommand(NewMessageCommand())
	rootCmd.AddCommand(NewRegistrationCommand())
}

func getHubNameAndConnectionString() (string, string, error) {
	hubName := viper.GetString("hub_name")
	if hubName == "" {
		return "", "", fmt.Errorf("hub name is required")
	}

	connectionString := viper.GetString("connection_string")
	if connectionString == "" {
		return "", "", fmt.Errorf("connection string is required")
	}

	return hubName, connectionString, nil
}

func getVerboseLogger() io.Writer {
	verbose := viper.GetBool("verbose")
	if verbose {
		return os.Stderr
	}
	return nil
}
