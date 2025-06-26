# Part 5: Configuration Management (The `config` command)

This section explains how `restapisummarizer` manages configuration, such as the Google Gemini API key. The `config` command provides a user-friendly way to interact with the application's settings.

## 1. The User Interface (`cmd/config.go`)

The `cmd/config.go` file defines the user-facing `config` command and its subcommands (`set`, `get`, and `view`).

**Key Responsibilities:**

*   **Command Structure:** It creates the main `configCmd` and attaches the `setCmd`, `getCmd`, and `viewCmd` as its children.
*   **User Input:** It handles the arguments and flags for each subcommand. For example, the `set` command requires a key and a value.
*   **Calling Core Logic:** It calls the functions in `internal/config/config.go` to perform the actual configuration management tasks.

**Code Snippet (`cmd/config.go`):**

```go
var setCmd = &cobra.Command{
    Use:   "set [key] [value]",
    Short: "Set a configuration value",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        key := args[0]
        value := args[1]

        if err := config.Set(key, value); err != nil {
            fmt.Printf("Error setting config: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Set %s = %s\n", key, value)
    },
}
```

## 2. The Engine (`internal/config/config.go`)

The heavy lifting of configuration management is done in `internal/config/config.go`. This file uses the `viper` library to handle reading from and writing to a configuration file.

**Key Responsibilities:**

*   **Viper Initialization:** It initializes `viper`, setting the configuration file path (`~/.restapisummarizer/config.yaml`) and enabling support for environment variables.
*   **Reading Configuration:** The `Get` function retrieves a configuration value. `viper` automatically handles looking for the value in the config file and then in environment variables.
*   **Writing Configuration:** The `Set` function writes a key-value pair to the configuration file.

**Key Features of `viper`:**

*   **Priority Order:** `viper` reads configuration from multiple sources with a defined priority: command-line flags, environment variables, configuration file, and default values.
*   **Automatic Loading:** It can automatically load configuration from a file if it exists.
*   **Environment Variable Binding:** It can automatically bind environment variables to configuration keys.

By separating the user-facing commands from the configuration management logic, the code is kept clean and organized. The use of `viper` provides a powerful and flexible way to handle configuration, making the application easy to configure for different environments.

