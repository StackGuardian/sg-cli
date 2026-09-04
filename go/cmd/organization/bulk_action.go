// Code generated for sg-cli; edit the generator spec or this file directly.

package organization

import (
	"context"
	"fmt"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/organizations"
	"github.com/spf13/cobra"
)

func newBulkActionCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-action",
		Short: "Activate or deactivate workflows and stacks in bulk",
		Long:  `Activate or deactivate workflows and stacks in bulk`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			action := common.String(cmd, "action")
			resourceIds := common.OptStringSlice(cmd, "resource-ids")
			if action != "activate" && action != "deactivate" {
				return fmt.Errorf("--action must be activate or deactivate, got %q", action)
			}
			response, err := c.Organizations.BulkActionWorkflows(ctx, org, &organizations.WorkflowsBulkActionRequest{Action: action, ResourceIds: resourceIds})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Bulk action submitted.")
			return nil
		},
	}
	cmd.Flags().String("org", "", "The organization name on Stackguardian platform.")
	_ = cmd.MarkFlagRequired("org")
	cmd.Flags().String("action", "", "activate or deactivate.")
	_ = cmd.MarkFlagRequired("action")
	cmd.Flags().StringSlice("resource-ids", nil, "Resource paths relative to the organization, e.g. /wfgrps/g/wfs/w (comma-separated).")
	_ = cmd.MarkFlagRequired("resource-ids")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
