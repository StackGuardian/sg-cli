// Code generated for sg-cli; edit the generator spec or this file directly.

package runnergroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/runnergroups"
	"github.com/spf13/cobra"
)

func newListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List runner groups",
		Long:  `List runner groups`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.RunnerGroups.ListAllRunnerGroups(ctx, org, &runnergroups.ListAllRunnerGroupsRequest{Limit: common.OptInt(cmd, "limit"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), FilterResourceNames: common.OptString(cmd, "resource-names")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("resource-names", "", "Filter by resource names.")
	return cmd
}
