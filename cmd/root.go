package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tarantino19/restgo/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "restapisummarizer",
	Short: "A CLI tool to analyze and summarize REST API endpoints",
	Long: `REST API Summarizer is a CLI tool that scans your codebase,
finds all REST API endpoints, and generates AI-powered summaries 
of what each endpoint does using Google's Gemini API.`,
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	
	// Add persistent flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.restapisummarizer/config.yaml)")
} 