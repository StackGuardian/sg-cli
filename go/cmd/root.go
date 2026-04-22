package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/StackGuardian/sg-cli/cmd/artifacts"
	"github.com/StackGuardian/sg-cli/cmd/gitscan"
	"github.com/StackGuardian/sg-cli/cmd/interactive"
	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-cli/cmd/stack"
	workflow "github.com/StackGuardian/sg-cli/cmd/workflow"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/option"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const version = "2.1.0"

var rootCmd = &cobra.Command{
	Use:     "sg-cli",
	Version: version,
	Short:   "Manage resources on the StackGuardian platform.",
	Long:    output.Banner(version),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	setHelpTemplate(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	API_KEY := "apikey " + os.Getenv("SG_API_TOKEN")
	SG_BASE_URL := os.Getenv("SG_BASE_URL")
	if SG_BASE_URL == "" {
		SG_BASE_URL = "https://api.app.stackguardian.io"
	}
	c := client.NewClient(
		option.WithApiKey(API_KEY),
		option.WithBaseURL(SG_BASE_URL),
	)
	rootCmd.AddCommand(workflow.NewWorkflowCmd(c))
	rootCmd.AddCommand(stack.NewStackCmd(c))
	rootCmd.AddCommand(artifacts.NewArtifactsCmd(c))
	rootCmd.AddCommand(gitscan.NewGitScanCmd())
	rootCmd.AddCommand(interactive.NewInteractiveCmd(c))
}

// setHelpTemplate applies a styled help template to cmd and all its descendants.
func setHelpTemplate(cmd *cobra.Command) {
	cmd.SetHelpFunc(styledHelp)
	for _, sub := range cmd.Commands() {
		setHelpTemplate(sub)
	}
}

func styledHelp(cmd *cobra.Command, args []string) {
	// Banner only on root
	if !cmd.HasParent() {
		fmt.Println(output.Banner(version))
		fmt.Println()
	}

	// Command path + short description
	path := output.PrimaryStyle.Render(cmd.CommandPath())
	if cmd.Short != "" {
		fmt.Println(path + output.Muted.Render("  —  ") + output.Muted.Render(cmd.Short))
	} else {
		fmt.Println(path)
	}
	fmt.Println()

	// Long description (skip on root since banner already shown)
	if cmd.Long != "" && cmd.HasParent() {
		fmt.Println(output.Muted.Render(cmd.Long))
		fmt.Println()
	}

	// Usage line
	if cmd.Runnable() || cmd.HasAvailableSubCommands() {
		fmt.Println(output.SectionStyle.Render("Usage"))
		fmt.Println()
		fmt.Println("  " + output.Bold.Render(cmd.UseLine()))
		if cmd.HasAvailableSubCommands() {
			fmt.Println("  " + output.Bold.Render(cmd.CommandPath()+" [command]"))
		}
		fmt.Println()
	}

	// Subcommands
	if cmd.HasAvailableSubCommands() {
		fmt.Println(output.SectionStyle.Render("Commands"))
		fmt.Println()
		for _, sub := range cmd.Commands() {
			if sub.IsAvailableCommand() {
				name := lipgloss.NewStyle().Foreground(lipgloss.Color("#88B7DA")).Bold(true).Width(16).Render(sub.Name())
				desc := output.Muted.Render(sub.Short)
				fmt.Println("  " + name + desc)
			}
		}
		fmt.Println()
	}

	// For commands that act as group parents (have subcommands + persistent flags like --org),
	// show persistent flags as "Required Flags". For leaf commands, show all local flags together.
	if cmd.HasAvailableSubCommands() && cmd.HasAvailablePersistentFlags() {
		fmt.Println(output.SectionStyle.Render("Required Flags"))
		fmt.Println()
		printFlags(cmd.PersistentFlags().FlagUsages())
		fmt.Println()
		// Also show any true local flags (like --help)
		localOnly := cmd.LocalNonPersistentFlags()
		if localOnly.HasAvailableFlags() {
			fmt.Println(output.SectionStyle.Render("Flags"))
			fmt.Println()
			printFlags(localOnly.FlagUsages())
			fmt.Println()
		}
	} else {
		// Leaf command: show all local flags
		if cmd.HasAvailableLocalFlags() {
			fmt.Println(output.SectionStyle.Render("Flags"))
			fmt.Println()
			printFlags(cmd.LocalFlags().FlagUsages())
			fmt.Println()
		}
		// Inherited global flags
		if cmd.HasAvailableInheritedFlags() {
			fmt.Println(output.SectionStyle.Render("Global Flags"))
			fmt.Println()
			printFlags(cmd.InheritedFlags().FlagUsages())
			fmt.Println()
		}
	}

	// Help hint
	if cmd.HasAvailableSubCommands() {
		fmt.Println(output.HintStyle.Render(
			fmt.Sprintf("Use \"%s [command] --help\" for more information about a command.", cmd.CommandPath()),
		))
		fmt.Println()
	}
}

func printFlags(usage string) {
	for _, line := range strings.Split(strings.TrimRight(usage, "\n"), "\n") {
		if line == "" {
			continue
		}
		// Colorize the flag name (the --flag part)
		if idx := strings.Index(line, "--"); idx >= 0 {
			before := line[:idx]
			rest := line[idx:]
			// Split flag name from the rest
			spaceIdx := strings.IndexAny(rest, " \t")
			if spaceIdx >= 0 {
				flagName := lipgloss.NewStyle().Foreground(lipgloss.Color("#88B7DA")).Render(rest[:spaceIdx])
				fmt.Println(before + flagName + output.Muted.Render(rest[spaceIdx:]))
				continue
			}
		}
		fmt.Println(output.Muted.Render(line))
	}
}
