package gemini

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/tarantino19/restgo/pkg/models"
	"google.golang.org/api/option"
)

// Client wraps the Gemini API client
type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewClient creates a new Gemini client
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Use Gemini 1.5 Flash for faster responses
	model := client.GenerativeModel("gemini-1.5-flash")
	
	// Configure model parameters for concise output
	model.SetTemperature(0.2) // Lower temperature for more consistent output
	model.SetMaxOutputTokens(100) // Limit output tokens per endpoint

	return &Client{
		client: client,
		model:  model,
	}, nil
}

// Close closes the Gemini client
func (c *Client) Close() error {
	return c.client.Close()
}

// SummarizeEndpoints processes multiple endpoints efficiently
func (c *Client) SummarizeEndpoints(ctx context.Context, endpoints []*models.Endpoint) error {
	// Group endpoints into batches for efficient processing
	batchSize := 5 // Process 5 endpoints at once
	batches := batchEndpoints(endpoints, batchSize)
	
	// Process batches concurrently with rate limiting
	semaphore := make(chan struct{}, 3) // Max 3 concurrent requests
	var wg sync.WaitGroup
	
	for i, batch := range batches {
		wg.Add(1)
		go func(batchIndex int, batchEndpoints []*models.Endpoint) {
			defer wg.Done()
			
			// Rate limiting: wait between batches
			if batchIndex > 0 {
				time.Sleep(2 * time.Second) // Reduced from 4 seconds
			}
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			// Process batch
			if err := c.summarizeBatch(ctx, batchEndpoints); err != nil {
				fmt.Printf("Warning: Failed to summarize batch %d: %v\n", batchIndex, err)
			}
		}(i, batch)
	}
	
	wg.Wait()
	return nil
}

// summarizeBatch processes a batch of endpoints in a single request
func (c *Client) summarizeBatch(ctx context.Context, endpoints []*models.Endpoint) error {
	if len(endpoints) == 0 {
		return nil
	}
	
	// Build optimized batch prompt
	prompt := buildBatchPrompt(endpoints)
	
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}
	
	// Parse batch response
	summaries := parseBatchResponse(resp)
	
	// Assign summaries to endpoints
	for i, endpoint := range endpoints {
		if i < len(summaries) {
			endpoint.Summary = summaries[i]
		} else {
			endpoint.Summary = "Summary unavailable"
		}
	}
	
	return nil
}

// buildBatchPrompt creates an optimized prompt for multiple endpoints
func buildBatchPrompt(endpoints []*models.Endpoint) string {
	var builder strings.Builder
	
	builder.WriteString("Analyze these REST API endpoints. For each, provide a one-line summary (max 50 chars).\n")
	builder.WriteString("Format: [N] Summary\n\n")
	
	for i, endpoint := range endpoints {
		// Include only essential information to reduce tokens
		builder.WriteString(fmt.Sprintf("[%d] %s %s\n", i+1, endpoint.Method, endpoint.Path))
		
		// Extract only the most relevant code lines
		relevantCode := extractRelevantCode(endpoint.RawCode)
		if relevantCode != "" {
			builder.WriteString(fmt.Sprintf("Code: %s\n", relevantCode))
		}
		builder.WriteString("\n")
	}
	
	builder.WriteString("Summaries:")
	
	return builder.String()
}

// extractRelevantCode extracts only the most relevant part of the code
func extractRelevantCode(rawCode string) string {
	lines := strings.Split(rawCode, "\n")
	var relevant []string
	
	// Look for function definitions, return statements, or key operations
	keywords := []string{"def ", "function", "return", "create", "update", "delete", "get", "post", "find", "save"}
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		
		// Check if line contains relevant keywords
		lower := strings.ToLower(trimmed)
		for _, keyword := range keywords {
			if strings.Contains(lower, keyword) {
				relevant = append(relevant, trimmed)
				if len(relevant) >= 3 { // Limit to 3 most relevant lines
					break
				}
			}
		}
	}
	
	// If no relevant lines found, take the first non-empty line
	if len(relevant) == 0 && len(lines) > 0 {
		for _, line := range lines {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				relevant = append(relevant, trimmed)
				break
			}
		}
	}
	
	return strings.Join(relevant, "; ")
}

// parseBatchResponse extracts summaries from batch response
func parseBatchResponse(resp *genai.GenerateContentResponse) []string {
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return []string{}
	}
	
	var content strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			content.WriteString(string(text))
		}
	}
	
	// Parse numbered summaries
	var summaries []string
	lines := strings.Split(content.String(), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for patterns like [1], [2], etc.
		if strings.HasPrefix(line, "[") {
			if idx := strings.Index(line, "]"); idx > 0 {
				summary := strings.TrimSpace(line[idx+1:])
				if summary != "" {
					// Ensure summary is not too long
					if len(summary) > 50 {
						summary = summary[:47] + "..."
					}
					summaries = append(summaries, summary)
				}
			}
		}
	}
	
	return summaries
}

// batchEndpoints groups endpoints into batches
func batchEndpoints(endpoints []*models.Endpoint, batchSize int) [][]*models.Endpoint {
	var batches [][]*models.Endpoint
	
	for i := 0; i < len(endpoints); i += batchSize {
		end := i + batchSize
		if end > len(endpoints) {
			end = len(endpoints)
		}
		batches = append(batches, endpoints[i:end])
	}
	
	return batches
} 