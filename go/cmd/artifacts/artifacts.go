package artifacts

import (
	"github.com/StackGuardian/sg-cli/cmd/artifacts/list"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewArtifactsCmd(c *client.Client) *cobra.Command {
	var artifactsCmd = &cobra.Command{
		Use:   "artifacts",
		Short: "Manage artifacts",
		Long:  "List and manage artifacts on the StackGuardian platform.",
	}

	artifactsCmd.PersistentFlags().String("org", "", "The organization name on StackGuardian platform.")
	artifactsCmd.MarkPersistentFlagRequired("org")

	artifactsCmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	artifactsCmd.MarkPersistentFlagRequired("workflow-group")

	artifactsCmd.PersistentFlags().String("workflow-id", "", "The workflow id in the workflow group.")
	artifactsCmd.MarkPersistentFlagRequired("workflow-id")

	artifactsCmd.AddCommand(list.NewListCmd(c))

	return artifactsCmd
}
