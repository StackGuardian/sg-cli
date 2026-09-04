// Code generated for sg-cli; edit the generator spec or this file directly.

package stackrun

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewStackRunCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack-run",
		Short: "Trigger and inspect stack runs",
		Long:  `Trigger and inspect stack runs`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Trigger a stack run (use --action-type or a payload)
  read                   Get a stack run
  list                   List the runs of a stack`)
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
		newListCmd(c),
	)
	return cmd
}
