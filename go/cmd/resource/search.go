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

func newSearchCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search -- payload.json",
		Short: "Search resources (payload and flags are merged; ResourceTypes is required by the API)",
		Long:  `Search resources (payload and flags are merged; ResourceTypes is required by the API)`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req resources.SearchResourcesRequest
			if len(args) > 0 {
				if err := common.ReadPayload(cmd, args, &req); err != nil {
					if errors.Is(err, common.ErrDryRun) {
						return nil
					}
					return err
				}
			}
			if v := common.OptString(cmd, "search-query"); v != nil {
				req.SearchQuery = v
			}
			if v := common.OptStringSlice(cmd, "resource-types"); v != nil {
				req.ResourceTypes = v
			}
			if v := common.OptInt(cmd, "limit"); v != nil {
				req.Limit = v
			}
			if v := common.OptString(cmd, "last-evaluated-key"); v != nil {
				req.LastEvaluatedKey = v
			}
			response, err := c.Resources.SearchResources(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Search completed.")
			return nil
		},
	}
	cmd.Flags().String("search-query", "", "Free-text search query.")
	cmd.Flags().StringSlice("resource-types", nil, "Resource types to search, e.g. WORKFLOW,STACK (comma-separated).")
	cmd.Flags().Int("limit", 0, "Maximum number of results.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
