package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tarantino19/restgo/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Manage REST API Summarizer configuration settings, including API keys.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
}

var setAPIKeyCmd = &cobra.Command{
	Use:   "api-key [key]",
	Short: "Set the Gemini API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := args[0]
		
		err := config.SetAPIKey(apiKey)
		if err != nil {
			color.Red("Error saving API key: %v", err)
			return
		}
		
		color.Green("✓ API key saved successfully!")
		color.Yellow("You can now use 'restapisummarizer sum' to analyze your APIs.")
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get configuration values",
}

var getAPIKeyCmd = &cobra.Command{
	Use:   "api-key",
	Short: "Get the current Gemini API key",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := config.GetAPIKey()
		
		if apiKey == "" {
			color.Yellow("No API key configured.")
			color.Yellow("Set one using: restapisummarizer config set api-key YOUR_KEY")
			return
		}
		
		// Mask the API key for security
		maskedKey := maskAPIKey(apiKey)
		fmt.Printf("Current API key: %s\n", maskedKey)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configSetCmd.AddCommand(setAPIKeyCmd)
	configGetCmd.AddCommand(getAPIKeyCmd)
}

// maskAPIKey masks an API key for display
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
} 