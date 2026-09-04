package cmd

import (
	"os"

	"github.com/StackGuardian/sg-cli/cmd/apiaccess"
	"github.com/StackGuardian/sg-cli/cmd/apitoken"
	"github.com/StackGuardian/sg-cli/cmd/artifacts"
	"github.com/StackGuardian/sg-cli/cmd/auditlog"
	"github.com/StackGuardian/sg-cli/cmd/benchmarkreport"
	"github.com/StackGuardian/sg-cli/cmd/billing"
	"github.com/StackGuardian/sg-cli/cmd/chat"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-cli/cmd/connector"
	"github.com/StackGuardian/sg-cli/cmd/connectorgroup"
	"github.com/StackGuardian/sg-cli/cmd/organization"
	"github.com/StackGuardian/sg-cli/cmd/policy"
	"github.com/StackGuardian/sg-cli/cmd/resource"
	"github.com/StackGuardian/sg-cli/cmd/role"
	"github.com/StackGuardian/sg-cli/cmd/rolebinding"
	"github.com/StackGuardian/sg-cli/cmd/runnergroup"
	"github.com/StackGuardian/sg-cli/cmd/secret"
	"github.com/StackGuardian/sg-cli/cmd/stack"
	"github.com/StackGuardian/sg-cli/cmd/stackrun"
	"github.com/StackGuardian/sg-cli/cmd/stacktemplate"
	"github.com/StackGuardian/sg-cli/cmd/stackworkflow"
	"github.com/StackGuardian/sg-cli/cmd/stackworkflowrun"
	"github.com/StackGuardian/sg-cli/cmd/statebackend"
	"github.com/StackGuardian/sg-cli/cmd/template"
	"github.com/StackGuardian/sg-cli/cmd/user"
	workflow "github.com/StackGuardian/sg-cli/cmd/workflow"
	"github.com/StackGuardian/sg-cli/cmd/workflowgroup"
	"github.com/StackGuardian/sg-cli/cmd/workflowrun"
	"github.com/StackGuardian/sg-cli/cmd/workflowsteptemplate"
	"github.com/StackGuardian/sg-cli/cmd/workflowtemplate"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "sg-cli",
	Version: "2.2.0",
	Short:   "sg-cli is CLI command for managing resources on Stackguardian platform.",
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
	c := common.NewClient()
	rootCmd.AddCommand(workflow.NewWorkflowCmd(c))
	rootCmd.AddCommand(stack.NewStackCmd(c))
	rootCmd.AddCommand(artifacts.NewArtifactsCmd(c))
	rootCmd.AddCommand(apiaccess.NewApiAccessCmd(c))
	rootCmd.AddCommand(apitoken.NewApiTokenCmd(c))
	rootCmd.AddCommand(auditlog.NewAuditLogCmd(c))
	rootCmd.AddCommand(benchmarkreport.NewBenchmarkReportCmd(c))
	rootCmd.AddCommand(billing.NewBillingCmd(c))
	rootCmd.AddCommand(chat.NewChatCmd(c))
	rootCmd.AddCommand(connector.NewConnectorCmd(c))
	rootCmd.AddCommand(connectorgroup.NewConnectorGroupCmd(c))
	rootCmd.AddCommand(organization.NewOrganizationCmd(c))
	rootCmd.AddCommand(policy.NewPolicyCmd(c))
	rootCmd.AddCommand(resource.NewResourceCmd(c))
	rootCmd.AddCommand(role.NewRoleCmd(c))
	rootCmd.AddCommand(rolebinding.NewRoleBindingCmd(c))
	rootCmd.AddCommand(runnergroup.NewRunnerGroupCmd(c))
	rootCmd.AddCommand(secret.NewSecretCmd(c))
	rootCmd.AddCommand(stackrun.NewStackRunCmd(c))
	rootCmd.AddCommand(stacktemplate.NewStackTemplateCmd(c))
	rootCmd.AddCommand(stackworkflow.NewStackWorkflowCmd(c))
	rootCmd.AddCommand(stackworkflowrun.NewStackWorkflowRunCmd(c))
	rootCmd.AddCommand(statebackend.NewStateBackendCmd(c))
	rootCmd.AddCommand(template.NewTemplateCmd(c))
	rootCmd.AddCommand(user.NewUserCmd(c))
	rootCmd.AddCommand(workflowgroup.NewWorkflowGroupCmd(c))
	rootCmd.AddCommand(workflowrun.NewWorkflowRunCmd(c))
	rootCmd.AddCommand(workflowsteptemplate.NewWorkflowStepTemplateCmd(c))
	rootCmd.AddCommand(workflowtemplate.NewWorkflowTemplateCmd(c))
}

// RootCommand returns the root command with every sub-command attached (used by tests).
func RootCommand() *cobra.Command {
	return rootCmd
}
