package analyzer

import (
	"bufio"
	"fmt"
	"os" // Keep os for os.ReadFile
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/tarantino19/restgo/pkg/models"
)

// Analyzer handles the analysis of source code files
type Analyzer struct {
	patterns []FrameworkPatterns
}

// NewAnalyzer creates a new analyzer instance
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		patterns: GetAllPatterns(),
	}
}

// AnalyzeDirectory scans a directory for REST API endpoints
func (a *Analyzer) AnalyzeDirectory(dir string) ([]*models.Endpoint, error) {
	var endpoints []*models.Endpoint
	var filesAnalyzed int
	var currentDir string

	color.Blue("🔍 Scanning directory tree: %s", dir)
	color.Blue("This will recursively scan all subdirectories...\n")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Track current directory for progress display
		dir := filepath.Dir(path)
		if dir != currentDir && info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			currentDir = dir
			relPath, _ := filepath.Rel(dir, path)
			if relPath != "." {
				color.HiBlack("  📂 Entering: %s", path)
			}
		}

		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Skip common non-source directories
		if shouldSkipPath(path) {
			return nil
		}

		// Check if file matches any framework patterns
		ext := filepath.Ext(path)
		for _, framework := range a.patterns {
			for _, filePattern := range framework.FilePatterns {
				if ext == filePattern {
					filesAnalyzed++
					fileEndpoints, err := a.analyzeFile(path, framework)
					if err != nil {
						color.Yellow("Warning: Error analyzing %s: %v", path, err)
						continue
					}
					endpoints = append(endpoints, fileEndpoints...)
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	color.Green("\n✓ Scan complete! Analyzed %d files, found %d endpoints", filesAnalyzed, len(endpoints))
	return endpoints, nil
}

// analyzeFile analyzes a single file for endpoints
func (a *Analyzer) analyzeFile(filePath string, framework FrameworkPatterns) ([]*models.Endpoint, error) {
	content, err := os.ReadFile(filePath) // Changed from ioutil.ReadFile
	if err != nil {
		return nil, err
	}

	var endpoints []*models.Endpoint
	lines := strings.Split(string(content), "\n")

	for lineNum, line := range lines {
		for _, pattern := range framework.Patterns {
			matches := pattern.Regex.FindStringSubmatch(line)
			if matches != nil {
				endpoint := a.extractEndpoint(matches, pattern, filePath, lineNum+1, lines, framework.Name)
				if endpoint != nil {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	if len(endpoints) > 0 {
		// Show relative path for better visibility of subdirectories
		cwd, _ := os.Getwd()
		relPath, _ := filepath.Rel(cwd, filePath)
		color.Cyan("  ✓ Found %d endpoints in %s", len(endpoints), relPath)
	}

	return endpoints, nil
}

// extractEndpoint extracts endpoint information from regex matches
func (a *Analyzer) extractEndpoint(matches []string, pattern Pattern, filePath string, lineNum int, lines []string, framework string) *models.Endpoint {
	var method, path string

	// Extract method and path based on pattern configuration
	if pattern.MethodIndex > 0 && pattern.MethodIndex < len(matches) {
		method = strings.ToUpper(matches[pattern.MethodIndex])
	}

	if pattern.PathIndex > 0 && pattern.PathIndex < len(matches) {
		path = matches[pattern.PathIndex]
	}

	// Handle special cases
	if pattern.MethodIndex == 0 && strings.Contains(lines[lineNum-1], "resources") {
		// Rails resources generates multiple endpoints
		return &models.Endpoint{
			Method:    "RESOURCE",
			Path:      "/" + path,
			File:      filePath,
			Line:      lineNum,
			Framework: framework,
			Language:  getLanguageFromExtension(filepath.Ext(filePath)),
			RawCode:   extractCodeContext(lines, lineNum-1, 3),
		}
	}

	// Default to GET if method not found (e.g., Flask without methods specified)
	if method == "" && path != "" {
		method = "GET"
	}

	// Skip if we couldn't extract both method and path
	if method == "" || path == "" {
		return nil
	}

	return &models.Endpoint{
		Method:    method,
		Path:      path,
		File:      filePath,
		Line:      lineNum,
		Framework: framework,
		Language:  getLanguageFromExtension(filepath.Ext(filePath)),
		RawCode:   extractCodeContext(lines, lineNum-1, 5),
	}
}

// extractCodeContext extracts surrounding code for context
func extractCodeContext(lines []string, centerLine int, contextSize int) string {
	start := centerLine - contextSize
	if start < 0 {
		start = 0
	}

	end := centerLine + contextSize
	if end >= len(lines) {
		end = len(lines) - 1
	}

	var context []string
	for i := start; i <= end; i++ {
		// Skip empty lines and comments to save tokens
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") {
			context = append(context, lines[i])
		}
	}

	return strings.Join(context, "\n")
}

// shouldSkipPath checks if a path should be skipped
func shouldSkipPath(path string) bool {
	skipDirs := []string{
		"node_modules",
		"vendor",
		".git",
		"dist",
		"build",
		"target",
		"__pycache__",
		".venv",
		"venv",
		"env",
		".idea",
		".vscode",
		"coverage",
		"test",
		"tests",
		"spec",
		"specs",
		".next",
		"out",
		"tmp",
		"temp",
		"cache",
		".cache",
		"logs",
		"docs",
		"documentation",
		"examples",
		"migrations",
		"public",
		"static",
		"assets",
		"bin",
		"obj",
	}

	// Check if file is too large (skip files over 1MB)
	if info, err := os.Stat(path); err == nil && info.Size() > 1024*1024 {
		return true
	}

	// Check each part of the path
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		for _, skip := range skipDirs {
			if part == skip {
				return true
			}
		}
		// Skip hidden directories
		if strings.HasPrefix(part, ".") && part != "." && part != ".." {
			return true
		}
	}

	// Skip minified files
	if strings.Contains(path, ".min.") {
		return true
	}

	return false
}

// getLanguageFromExtension returns the language based on file extension
func getLanguageFromExtension(ext string) string {
	languages := map[string]string{
		".js":   "JavaScript",
		".ts":   "TypeScript",
		".mjs":  "JavaScript",
		".py":   "Python",
		".java": "Java",
		".go":   "Go",
		".rb":   "Ruby",
		".cs":   "C#",
		".php":  "PHP",
	}

	if lang, ok := languages[ext]; ok {
		return lang
	}

	return "Unknown"
}

// ReadFileLines reads a file and returns its lines for detailed analysis
func ReadFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
