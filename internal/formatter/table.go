package formatter

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/tarantino19/restgo/pkg/models"
)

// FormatEndpointsTable formats endpoints as a table
func FormatEndpointsTable(endpoints []*models.Endpoint) {
	if len(endpoints) == 0 {
		color.Yellow("No endpoints found.")
		return
	}

	// Print summary header
	color.Green("\n🔍 REST API Endpoints Summary\n")
	fmt.Printf("Found %d endpoints\n\n", len(endpoints))

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	
	// Add header
	table.Append([]string{
		color.CyanString("Method"),
		color.CyanString("Path"),
		color.CyanString("File"),
		color.CyanString("Summary"),
	})

	// Add data rows
	for _, endpoint := range endpoints {
		method := colorizeMethodSimple(endpoint.Method)
		path := endpoint.Path
		file := fmt.Sprintf("%s:%d", shortenPath(endpoint.File), endpoint.Line)
		summary := endpoint.Summary
		
		table.Append([]string{method, path, file, summary})
	}

	// Render table
	table.Render()
	
	// Print grouped by file
	printGroupedByFile(endpoints)
}

// shortenPath shortens file path for display
func shortenPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return ".../" + strings.Join(parts[len(parts)-2:], "/")
	}
	return path
}

// printGroupedByFile prints endpoints grouped by file
func printGroupedByFile(endpoints []*models.Endpoint) {
	color.Green("\n📁 Endpoints by File:\n")
	
	// Group by file
	fileMap := make(map[string][]*models.Endpoint)
	for _, endpoint := range endpoints {
		fileMap[endpoint.File] = append(fileMap[endpoint.File], endpoint)
	}
	
	// Print each file's endpoints
	for file, eps := range fileMap {
		color.Cyan("  %s (%d endpoints)\n", file, len(eps))
		for _, ep := range eps {
			fmt.Printf("    • %s %s\n", colorizeMethodSimple(ep.Method), ep.Path)
		}
	}
}

// colorizeMethodSimple returns a simple colored method string
func colorizeMethodSimple(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return color.GreenString(method)
	case "POST":
		return color.BlueString(method)
	case "PUT":
		return color.YellowString(method)
	case "DELETE":
		return color.RedString(method)
	case "PATCH":
		return color.MagentaString(method)
	default:
		return method
	}
} 