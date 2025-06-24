package models

// Endpoint represents a REST API endpoint
type Endpoint struct {
	Method      string // HTTP method (GET, POST, PUT, DELETE, etc.)
	Path        string // Endpoint path (e.g., /users/:id)
	File        string // Source file where endpoint is defined
	Line        int    // Line number in source file
	Function    string // Function/handler name
	Summary     string // AI-generated summary
	Language    string // Programming language
	Framework   string // Web framework used
	RawCode     string // Raw code snippet for context
} 