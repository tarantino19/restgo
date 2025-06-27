package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tarantino19/restgo/internal/analyzer"
	"github.com/tarantino19/restgo/internal/cache"
	"github.com/tarantino19/restgo/internal/config"
	"github.com/tarantino19/restgo/internal/formatter"
	"github.com/tarantino19/restgo/internal/gemini"
	"github.com/tarantino19/restgo/pkg/models"
)

var (
	noCache bool
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
	rootCmd.AddCommand(sumCmd) //viper native command AddComand
	sumCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable cache and regenerate all summaries")
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

	// Initialize cache
	endpointCache, err := cache.NewCache(24 * time.Hour) // Added 24 * time.Hour as expiration
	if err != nil {
		color.Yellow("Warning: Could not initialize cache: %v", err)
		// Continue without cache
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

	// Check cache for existing summaries
	var needsSummary []*models.Endpoint
	cachedCount := 0

	if endpointCache != nil && !noCache {
		for _, endpoint := range endpoints {
			fileHash := cache.HashFile(endpoint.RawCode)
			if summary, found := endpointCache.Get(endpoint.Method, endpoint.Path, fileHash); found {
				endpoint.Summary = summary
				cachedCount++
			} else {
				needsSummary = append(needsSummary, endpoint)
			}
		}

		if cachedCount > 0 {
			color.Blue("📦 Using %d cached summaries\n", cachedCount)
		}
	} else {
		needsSummary = endpoints
	}

	// Generate summaries only for endpoints that need them
	if len(needsSummary) > 0 {
		// Create Gemini client
		color.Blue("🤖 Initializing Gemini AI...\n")
		geminiClient, err := gemini.NewClient(apiKey)
		if err != nil {
			color.Red("Error creating Gemini client: %v", err)
			os.Exit(1)
		}
		defer geminiClient.Close()

		// Generate summaries
		color.Blue("📝 Generating summaries for %d new endpoints...\n", len(needsSummary))
		ctx := context.Background()

		// Add timeout for the entire operation
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute) // Reduced from 10 minutes
		defer cancel()

		err = geminiClient.SummarizeEndpoints(ctx, needsSummary)
		if err != nil {
			color.Red("Error generating summaries: %v", err)
			// Continue anyway - we can still show endpoints without summaries
		}

		// Cache the new summaries
		if endpointCache != nil && !noCache {
			for _, endpoint := range needsSummary {
				if endpoint.Summary != "" && endpoint.Summary != "Summary unavailable" {
					fileHash := cache.HashFile(endpoint.RawCode)
					endpointCache.Set(endpoint.Method, endpoint.Path, fileHash, endpoint.Summary)
				}
			}
		}
	}

	// Display results
	formatter.FormatEndpointsTable(endpoints)

	// Show statistics
	duration := time.Since(startTime)
	color.Green("\n✅ Analysis completed in %s\n", duration.Round(time.Second))

	if cachedCount > 0 {
		percentage := (cachedCount * 100) / len(endpoints)
		color.HiBlack("   • %d/%d summaries from cache (%d%%)", cachedCount, len(endpoints), percentage)
	}

	if len(needsSummary) > 0 {
		tokensEstimate := len(needsSummary) * 50 // Rough estimate
		color.HiBlack("   • Generated %d new summaries (~%d tokens used)", len(needsSummary), tokensEstimate)
	}
}
