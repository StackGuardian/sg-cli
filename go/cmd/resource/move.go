// Code generated for sg-cli; edit the generator spec or this file directly.

package resource

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/resources"
	"github.com/spf13/cobra"
)

func newMoveCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "move -- payload.json",
		Short: "Move workflows, stacks or workflow groups to another parent",
		Long:  `Move workflows, stacks or workflow groups to another parent`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req resources.MoveResourcesRequest
			if len(args) > 0 {
				if err := common.ReadPayload(cmd, args, &req); err != nil {
					if errors.Is(err, common.ErrDryRun) {
						return nil
					}
					return err
				}
			}
			if v := common.OptStringSlice(cmd, "source-ids"); v != nil {
				req.SourceResourceIds = v
			}
			if v := common.OptString(cmd, "destination"); v != nil {
				req.DestinationParentId = *v
			}
			if len(req.SourceResourceIds) == 0 || req.DestinationParentId == "" {
				return errors.New("provide --source-ids and --destination, or a payload with SourceResourceIds and DestinationParentId")
			}
			response, err := c.Resources.MoveResources(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Resources moved.")
			return nil
		},
	}
	cmd.Flags().StringSlice("source-ids", nil, "Resource ids to move (comma-separated).")
	cmd.Flags().String("destination", "", "Id of the destination parent.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
