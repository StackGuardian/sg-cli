package cmd

import (
	"os"

	"github.com/StackGuardian/sg-cli/cmd/artifacts"
	"github.com/StackGuardian/sg-cli/cmd/stack"
	workflow "github.com/StackGuardian/sg-cli/cmd/workflow"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/option"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sg-cli",
	Short: "sg-cli is CLI command for managing resources on Stackguardian platform.",
	Long: `sg-cli is CLI command for managing resources on Stackguardian platform.
More information available at: https://docs.qa.stackguardian.io/docs/`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
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
}
