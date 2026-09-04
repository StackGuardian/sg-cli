// Code generated for sg-cli; edit the generator spec or this file directly.

package stackworkflow

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewStackWorkflowCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack-workflow",
		Short: "Manage the workflows that belong to a stack",
		Long:  `Manage the workflows that belong to a stack`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a workflow inside a stack
  read                   Get a stack workflow
  update                 Update a stack workflow
  delete                 Delete a stack workflow
  list                   List the workflows of a stack
  outputs                Get the outputs of a stack workflow
  tfstate-upload-url     Get a signed URL to upload a Terraform state file
  file-upload-url        Get a signed URL to upload a file to a stack workflow
  compare                Compare a stack workflow against a template without applying changes
  vcs-triggers           Create VCS webhook triggers for a stack workflow
  artifacts              Manage the artifacts of a stack workflow`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.PersistentFlags().String("workflow-group", "", "The workflow group under the organization.")
	_ = cmd.MarkPersistentFlagRequired("workflow-group")
	cmd.PersistentFlags().String("stack-id", "", "The stack ID.")
	_ = cmd.MarkPersistentFlagRequired("stack-id")
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newListCmd(c),
		newOutputsCmd(c),
		newTfstateUploadUrlCmd(c),
		newFileUploadUrlCmd(c),
		newCompareCmd(c),
		newVcsTriggersCmd(c),
		newArtifactsCmd(c),
	)
	return cmd
}
func newArtifactsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifacts",
		Short: "Manage the artifacts of a stack workflow",
		Long:  `Manage the artifacts of a stack workflow`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  list                   List the artifacts of a stack workflow
  get-url                Get a signed download URL for an artifact
  read                   Read an artifact's content (or a signed URL on private runners)
  create                 Create or replace an artifact with the JSON payload as its content
  delete                 Delete an artifact
  lock                   Lock the stack workflow on an artifact
  unlock                 Release the stack workflow's lock on an artifact`)
		},
	}
	cmd.PersistentFlags().String("workflow-id", "", "The workflow ID.")
	_ = cmd.MarkPersistentFlagRequired("workflow-id")
	cmd.AddCommand(
		newArtifactsListCmd(c),
		newArtifactsGetUrlCmd(c),
		newArtifactsReadCmd(c),
		newArtifactsCreateCmd(c),
		newArtifactsDeleteCmd(c),
		newArtifactsLockCmd(c),
		newArtifactsUnlockCmd(c),
	)
	return cmd
}
