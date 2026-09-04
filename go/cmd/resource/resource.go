// Code generated for sg-cli; edit the generator spec or this file directly.

package resource

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewResourceCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Move and search resources across the organization",
		Long:  `Move and search resources across the organization`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  move                   Move workflows, stacks or workflow groups to another parent
  search                 Search resources (payload and flags are merged; ResourceTypes is required by the API)
  tags                   List the distinct tags used on resources`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newMoveCmd(c),
		newSearchCmd(c),
		newTagsCmd(c),
	)
	return cmd
}
