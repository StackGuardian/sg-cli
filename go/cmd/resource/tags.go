// Code generated for sg-cli; edit the generator spec or this file directly.

package resource

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/resources"
	"github.com/spf13/cobra"
)

func newTagsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List the distinct tags used on resources",
		Long:  `List the distinct tags used on resources`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.Resources.ListResourceTags(ctx, org, &resources.ListResourceTagsRequest{ContextTagKey: common.OptString(cmd, "context-tag-key"), Limit: common.OptInt(cmd, "limit"), Query: common.OptString(cmd, "query"), ResourceType: common.OptString(cmd, "resource-type"), TagKind: common.OptString(cmd, "tag-kind")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("context-tag-key", "", "Context tag key to list values for.")
	cmd.Flags().Int("limit", 0, "Maximum number of tags.")
	cmd.Flags().String("query", "", "Filter tags by prefix or substring.")
	cmd.Flags().String("resource-type", "", "Restrict to a resource type.")
	cmd.Flags().String("tag-kind", "", "Tag kind to list.")
	return cmd
}
