package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tarantino19/restgo/internal/analyzer"
	"github.com/tarantino19/restgo/internal/config"
	"github.com/tarantino19/restgo/internal/formatter"
	"github.com/tarantino19/restgo/internal/gemini"
)

var sumCmd = &cobra.Command{
	Use:   "sum [directory]",
	Short: "Analyze REST API endpoints in a directory",
	Long: `Scans the specified directory (or current directory if not specified) 
for REST API endpoints and generates AI-powered summaries using Gemini API.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runSum,
}

func init() {
	rootCmd.AddCommand(sumCmd)
}

func runSum(cmd *cobra.Command, args []string) {
	// Determine directory to analyze
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Convert to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		color.Red("Error resolving directory path: %v", err)
		os.Exit(1)
	}

	// Check if directory exists
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		color.Red("Directory does not exist: %s", absDir)
		os.Exit(1)
	}

	// Check for API key
	apiKey := config.GetAPIKey()
	if apiKey == "" {
		color.Red("Error: Gemini API key not set")
		color.Yellow("Please set your API key using one of the following methods:")
		color.Yellow("1. Run: restapisummarizer config set api-key YOUR_API_KEY")
		color.Yellow("2. Set environment variable: export GEMINI_API_KEY=YOUR_API_KEY")
		os.Exit(1)
	}

	// Create analyzer
	analyzer := analyzer.NewAnalyzer()
	
	// Analyze directory
	color.Green("\n🚀 Starting REST API analysis...\n")
	startTime := time.Now()
	
	endpoints, err := analyzer.AnalyzeDirectory(absDir)
	if err != nil {
		color.Red("Error analyzing directory: %v", err)
		os.Exit(1)
	}

	if len(endpoints) == 0 {
		color.Yellow("No REST API endpoints found in %s", absDir)
		color.Yellow("Make sure the directory contains source code with REST API definitions.")
		return
	}

	// Create Gemini client
	color.Blue("\n🤖 Initializing Gemini AI...\n")
	geminiClient, err := gemini.NewClient(apiKey)
	if err != nil {
		color.Red("Error creating Gemini client: %v", err)
		os.Exit(1)
	}
	defer geminiClient.Close()

	// Generate summaries
	color.Blue("📝 Generating endpoint summaries...\n")
	ctx := context.Background()
	
	// Add timeout for the entire operation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	
	err = geminiClient.SummarizeEndpoints(ctx, endpoints)
	if err != nil {
		color.Red("Error generating summaries: %v", err)
		// Continue anyway - we can still show endpoints without summaries
	}

	// Display results
	formatter.FormatEndpointsTable(endpoints)
	
	// Show timing
	duration := time.Since(startTime)
	color.Green("\n✅ Analysis completed in %s\n", duration.Round(time.Second))
} 