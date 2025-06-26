# Study Guide: Understanding the `restapisummarizer` Codebase

Welcome! This guide is designed to help you, a junior engineer, understand the `restapisummarizer` project step-by-step. We'll go from a high-level overview right down to the core logic.

Our goal is not to understand every single line of code at once, but to build a mental model of how the application works, where to find things, and how the different pieces fit together.

**Key Technologies:**
*   **Go:** The programming language used.
*   **Cobra:** A popular library for creating powerful and modern CLI applications in Go.
*   **Viper:** A companion to Cobra for handling configuration (from files, environment variables, etc.).
*   **Google Gemini API:** The AI service used to generate endpoint summaries.

---

## The Plan

Here is the roadmap we will follow. We'll go through each step, and you can ask questions at any point.

### **Part 1: The Big Picture (What does this thing do?)**

*   **Goal:** Understand the application from a user's perspective.
*   **Activities:**
    1.  Review the `README.md` to understand the project's purpose, features, and basic usage.
    2.  Run `restapisummarizer --help` to see all the available commands and flags.
    3.  Briefly discuss the main commands: `sum` and `config`.

### **Part 2: Project Structure (Where is everything?)**

*   **Goal:** Learn the layout of the codebase.
*   **Activities:**
    1.  Walk through the main directories (`/cmd`, `/internal`, `/pkg`, `/scripts`).
    2.  Explain the purpose of each directory and the Go convention of `internal` vs. `pkg`.
    3.  Look at `go.mod` to see our project's dependencies.
    4.  Identify the main entry point of the entire application: `main.go`.

### **Part 3: The Entry Point (How does it start?)**

*   **Goal:** Trace the first few steps of the application's execution.
*   **Activities:**
    1.  Start at `main.go` and see how it calls `cmd.Execute()`.
    2.  Move to `cmd/root.go` to understand how the main `restapisummarizer` command is created using `cobra`.
    3.  Examine `cmd/sum.go` and `cmd/config.go` to see how subcommands are defined and attached to the root command.

### **Part 4: The Core Logic (Following the `sum` command)**

*   **Goal:** Deep dive into the primary feature: summarizing API endpoints.
*   **Activities:**
    1.  **The Commander (`cmd/sum.go`):** See how the `sum` command is executed.
    2.  **The Orchestrator (`internal/analyzer/analyzer.go`):** Follow the logic into the `internal` package, where the main work begins. See how it scans directories and files.
    3.  **The Brain (`internal/analyzer/patterns.go`):** Understand how we use Regular Expressions (Regex) to find API endpoints for different frameworks (like Express, Gin, etc.).
    4.  **The AI Client (`internal/gemini/client.go`):** Learn how the application communicates with the Google Gemini API to get the AI-powered summaries.
    5.  **The Presenter (`internal/formatter/table.go`):** See how the final results are formatted into a clean table for the user.

### **Part 5: Configuration Management (The `config` command)**

*   **Goal:** Understand how the application manages settings, like the API key.
*   **Activities:**
    1.  Look at `cmd/config.go` to see the user-facing commands (`set`, `get`).
    2.  Dive into `internal/config/config.go` to understand how `viper` is used to load configuration from files (`~/.restapisummarizer/config.yaml`) and environment variables.

---

Let me know which part of the study guide you'd like to start with!
