// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflowrun

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewStackWorkflowRunCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack-workflow-run",
		Short: "Trigger, inspect and approve runs of a stack workflow",
		Long:  `Trigger, inspect and approve runs of a stack workflow`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Trigger a run of a stack workflow (payload optional)
  read                   Get a stack workflow run
  update                 Update a stack workflow run
  delete                 Cancel a stack workflow run through the DELETE endpoint
  logs                   Get the logs of a stack workflow run
  approve                Approve (or, with --reject, reject) a stack workflow run waiting for approval
  list                   List the runs of a stack workflow
  facts                  Read and write the facts of a stack workflow run`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	_ = cmd.MarkPersistentFlagRequired("workflow-group")
	cmd.PersistentFlags().String("stack-id", "", "The stack ID.")
	_ = cmd.MarkPersistentFlagRequired("stack-id")
	cmd.PersistentFlags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkPersistentFlagRequired("workflow-id")
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newLogsCmd(c),
		newApproveCmd(c),
		newListCmd(c),
		newFactsCmd(c),
	)
	return cmd
}
func newFactsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "facts",
		Short: "Read and write the facts of a stack workflow run",
		Long:  `Read and write the facts of a stack workflow run`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get the facts of a stack workflow run
  create                 Create or replace the facts of a stack workflow run
  update                 Update the facts of a stack workflow run`)
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
