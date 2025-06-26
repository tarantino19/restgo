# Part 4: The Core Logic (Following the `sum` command)

This part of the guide dives into the main functionality of `restapisummarizer`: summarizing API endpoints. We'll trace the execution flow from the moment a user runs the `sum` command to the point where a summary table is displayed.

## 1. The Commander (`cmd/sum.go`)

The journey begins in `cmd/sum.go`. This file, using the `cobra` library, defines the `sum` command, its flags, and what happens when it's executed.

**Key Responsibilities:**

*   **Command Definition:** It creates the `sumCmd` object, which represents the `sum` command.
*   **Flag Handling:** It defines flags like `--path`, `--language`, and `--output` that allow users to customize the command's behavior.
*   **Execution Logic:** The `Run` function of the `cobra.Command` is where the action starts. It reads the user's input (flags) and calls the core logic in the `internal/analyzer` package.

**Code Snippet (`cmd/sum.go`):**

```go
var sumCmd = &cobra.Command{
    Use:   "sum [flags]",
    Short: "Summarize REST API endpoints in a directory",
    Long:  `sum scans a directory for REST API endpoints in supported languages and frameworks, then generates a summary of each endpoint.`,
    Run: func(cmd *cobra.Command, args []string) {
        // ... (Flag handling and validation)

        // Create a new analyzer
        a := analyzer.NewAnalyzer(language, directory)

        // Run the analysis
        endpoints, err := a.Analyze()
        if err != nil {
            // ... (Error handling)
        }

        // ... (Formatting and output)
    },
}
```

## 2. The Orchestrator (`internal/analyzer/analyzer.go`)

Once the `sum` command is executed, it passes control to the `analyzer`. This is the central component that orchestrates the process of finding and summarizing endpoints.

**Key Responsibilities:**

*   **Directory Scanning:** It walks through the specified directory to find all the relevant source code files.
*   **File Analysis:** For each file, it reads the content and uses regular expressions to identify potential API endpoints.
*   **Endpoint Creation:** When an endpoint is found, it creates an `Endpoint` object (defined in `pkg/models/endpoint.go`) to store information like the file path, line number, HTTP method, and path.

**Code Snippet (`internal/analyzer/analyzer.go`):**

```go
// Analyze starts the analysis of the specified directory.
func (a *Analyzer) Analyze() ([]models.Endpoint, error) {
    var endpoints []models.Endpoint

    filepath.Walk(a.Directory, func(path string, info os.FileInfo, err error) error {
        // ... (File filtering logic)

        // Read the file content
        content, err := ioutil.ReadFile(path)
        if err != nil {
            return nil
        }

        // Find endpoints in the content
        matches := a.findEndpoints(string(content))

        // ... (Create Endpoint objects)
    })

    return endpoints, nil
}
```

## 3. The Brain (`internal/analyzer/patterns.go`)

The `analyzer` needs to know what to look for. That's where `patterns.go` comes in. This file contains the regular expressions (regex) used to detect API endpoints in different programming languages and frameworks.

**Key Responsibilities:**

*   **Regex Definitions:** It defines a map of regular expressions, where each regex is designed to match the syntax of a specific framework (e.g., Express.js, Gin, Flask).
*   **Pattern Selection:** The `analyzer` selects the appropriate regex based on the language specified by the user.

**Code Snippet (`internal/analyzer/patterns.go`):**

```go
var frameworkPatterns = map[string]string{
    "go-gin":      `r\.(GET|POST|PUT|DELETE|PATCH)\("([^"]+)",`,
    "js-express":  `app\.(get|post|put|delete|patch)\(['"]([^'"]+)['"],`,
    // ... (Other frameworks)
}
```

## 4. The AI Client (`internal/gemini/client.go`)

After the `analyzer` has identified all the endpoints, it's time to generate the summaries. This is handled by the `gemini` client, which communicates with the Google Gemini API.

**Key Responsibilities:**

*   **API Communication:** It sends the code snippet containing the endpoint definition to the Gemini API.
*   **Prompt Engineering:** It constructs a specific prompt that asks the AI to summarize the endpoint's purpose.
*   **Summary Retrieval:** It receives the summary from the API and adds it to the `Endpoint` object.

## 5. The Presenter (`internal/formatter/table.go`)

Finally, with all the data collected and summarized, the results need to be presented to the user. The `table.go` file is responsible for formatting the endpoints into a clean, readable table.

**Key Responsibilities:**

*   **Table Formatting:** It uses a table-writer library to create a well-structured table with headers for the file, line number, method, path, and summary.
*   **Output Rendering:** It prints the formatted table to the console.

This completes the journey of the `sum` command. By understanding these five key components, you have a solid grasp of the core functionality of `restapisummarizer`.
