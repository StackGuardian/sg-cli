package artifacts

import (
	"fmt"

	"github.com/StackGuardian/sg-cli/cmd/artifacts/list"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewArtifactsCmd(c *client.Client) *cobra.Command {
	// artifactsCmd represents the Artifacts command
	var artifactsCmd = &cobra.Command{
		Use:   "artifacts",
		Short: "List Artifacts",
		Long:  `List artifacts on Stackguardian platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`Sub-commands:
  list      List Artifacts`)
		},
	}

	artifactsCmd.PersistentFlags().String("org", "", "The organization name on Stackguardian platform.")
	artifactsCmd.MarkPersistentFlagRequired("org")

	artifactsCmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	artifactsCmd.MarkPersistentFlagRequired("workflow-group")

	artifactsCmd.PersistentFlags().String("workflow-id", "", "The workflow id in the workflow group.")
	artifactsCmd.MarkPersistentFlagRequired("workflow-id")

	artifactsCmd.AddCommand(list.NewListCmd(c))

	return artifactsCmd
}
