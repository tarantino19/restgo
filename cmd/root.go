package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tarantino19/restgo/internal/config"
)

const customHelpTemplate = `{{.Long}}

{{if or .Runnable .HasAvailableSubCommands}}{{bold "Usage:"}}
  {{if .Runnable}}{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}{{.CommandPath}} [command]{{end}}{{end}}

{{if .HasAvailableSubCommands}}{{bold "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{cyan (rpad .Name .NamePadding) }} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableLocalFlags}}{{bold "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasAvailableInheritedFlags}}{{bold "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasHelpSubCommands}}{{bold "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var rootCmd = &cobra.Command{
	Use:   "restapisummarizer",
	Short: "A CLI tool to analyze and summarize REST API endpoints",
	Long: color.New(color.FgCyan, color.Bold).Sprint("REST API Summarizer") + ` is a CLI tool that scans your codebase,
finds all REST API endpoints, and generates AI-powered summaries
of what each endpoint does using Google's Gemini API.`,
	Version: "1.0.1",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	//initialize init config
	cobra.OnInitialize(config.InitConfig)

	cobra.AddTemplateFunc("bold", color.New(color.Bold).Sprint)
	cobra.AddTemplateFunc("cyan", color.New(color.FgCyan).Sprint)
	rootCmd.SetHelpTemplate(customHelpTemplate)

	// Add persistent flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.restapisummarizer/config.yaml)")
}
