package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Config struct {
	GeminiAPIKey string `mapstructure:"gemini_api_key"`
}

var (
	cfgFile string
	cfg     *Config
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".restapisummarizer".
		configPath := filepath.Join(home, ".restapisummarizer")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			os.Mkdir(configPath, 0755)
		}

		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("RESTAPI")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		color.Green("Using config file: %s", viper.ConfigFileUsed())
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
	}
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	if cfg == nil {
		InitConfig()
	}
	return cfg
}

// SetAPIKey saves the API key to the config file
func SetAPIKey(apiKey string) error {
	viper.Set("gemini_api_key", apiKey)

	// Ensure config directory exists
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".restapisummarizer")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.Mkdir(configPath, 0755); err != nil {
			return err
		}
	}

	configFile := filepath.Join(configPath, "config.yaml")
	return viper.WriteConfigAs(configFile)
}

// GetAPIKey returns the configured API key
func GetAPIKey() string {
	// First check environment variable
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		return apiKey
	}

	// Then check config
	if cfg != nil && cfg.GeminiAPIKey != "" {
		return cfg.GeminiAPIKey
	}

	// Finally check viper directly
	return viper.GetString("gemini_api_key")
}
