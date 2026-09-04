// Code generated for sg-cli; edit the generator spec or this file directly.

package policy

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewPolicyCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage policies",
		Long:  `Manage policies`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  create                 Create a new policy (payload must include PolicyType: GENERAL or FILTER.INSIGHT)
  read                   Get details of a policy
  update                 Update a policy (payload must include PolicyType)
  delete                 Delete a policy
  list                   List policies`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newCreateCmd(c),
		newReadCmd(c),
		newUpdateCmd(c),
		newDeleteCmd(c),
		newListCmd(c),
	)
	return cmd
}
