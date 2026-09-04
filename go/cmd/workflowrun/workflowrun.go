// Code generated for sg-cli; edit the generator spec or this file directly.

package workflowrun

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewWorkflowRunCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow-run",
		Short: "Trigger, inspect, cancel and approve workflow runs",
		Long:  `Trigger, inspect, cancel and approve workflow runs`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Trigger a workflow run (payload optional)
  read                   Get a workflow run
  update                 Update a workflow run
  cancel                 Cancel a workflow run
  delete                 Cancel a workflow run through the DELETE endpoint
  logs                   Get the logs of a workflow run
  approve                Approve (or, with --reject, reject) a workflow run waiting for approval
  list                   List the runs of a workflow
  read-by-ksuid          Get a workflow run by KSUIDs (runner and machine tokens only)
  update-by-ksuid        Update a workflow run by KSUIDs (runner and machine tokens only)
  facts                  Read and write the facts of a workflow run`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	cmd.PersistentFlags().String("workflow-id", "", "The workflow ID.")
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newCancelCmd(c),
		newDeleteCmd(c),
		newLogsCmd(c),
		newApproveCmd(c),
		newListCmd(c),
		newReadByKsuidCmd(c),
		newUpdateByKsuidCmd(c),
		newFactsCmd(c),
	)
	return cmd
}
func newFactsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facts",
		Short: "Read and write the facts of a workflow run",
		Long:  `Read and write the facts of a workflow run`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get the facts of a workflow run
  create                 Create or replace the facts of a workflow run
  update                 Update the facts of a workflow run`)
		},
	}
	cmd.PersistentFlags().String("run-id", "", "The workflow run ID.")
	_ = cmd.MarkPersistentFlagRequired("run-id")
	cmd.AddCommand(
		newFactsReadCmd(c),
		newFactsCreateCmd(c),
		newFactsUpdateCmd(c),
	)
	return cmd
}
