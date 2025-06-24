package gemini

import (
	"context"
	"fmt"
	"strings"
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
	
	// Configure model parameters
	model.SetTemperature(0.3)
	model.SetMaxOutputTokens(200)

	return &Client{
		client: client,
		model:  model,
	}, nil
}

// Close closes the Gemini client
func (c *Client) Close() error {
	return c.client.Close()
}

// SummarizeEndpoint generates a summary for a single endpoint
func (c *Client) SummarizeEndpoint(ctx context.Context, endpoint *models.Endpoint) error {
	prompt := buildPrompt(endpoint)
	
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	summary := extractSummary(resp)
	endpoint.Summary = summary
	
	return nil
}

// SummarizeEndpoints processes multiple endpoints with rate limiting
func (c *Client) SummarizeEndpoints(ctx context.Context, endpoints []*models.Endpoint) error {
	// Rate limiting: 15 requests per minute for free tier
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for i, endpoint := range endpoints {
		if i > 0 {
			select {
			case <-ticker.C:
				// Continue after rate limit delay
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err := c.SummarizeEndpoint(ctx, endpoint); err != nil {
			// Log error but continue with other endpoints
			fmt.Printf("Warning: Failed to summarize endpoint %s %s: %v\n", 
				endpoint.Method, endpoint.Path, err)
			endpoint.Summary = "Failed to generate summary"
		}
	}

	return nil
}

// buildPrompt creates the prompt for the Gemini API
func buildPrompt(endpoint *models.Endpoint) string {
	return fmt.Sprintf(`Analyze this REST API endpoint and provide a concise summary of what it does.

Framework: %s
HTTP Method: %s
Path: %s
Function/Handler: %s

Code:
%s

Provide a brief, one-line summary (max 100 characters) of what this endpoint does. Focus on the business logic, not technical implementation details. Be concise and clear.

Summary:`, 
		endpoint.Framework,
		endpoint.Method,
		endpoint.Path,
		endpoint.Function,
		endpoint.RawCode,
	)
}

// extractSummary extracts the summary from the Gemini response
func extractSummary(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "No summary available"
	}

	var summary strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			summary.WriteString(string(text))
		}
	}

	result := strings.TrimSpace(summary.String())
	
	// Ensure summary is not too long
	if len(result) > 100 {
		result = result[:97] + "..."
	}
	
	return result
} 