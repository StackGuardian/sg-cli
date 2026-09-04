// Code generated for sg-cli; edit the generator spec or this file directly.

package organization

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/organizations"
	"github.com/spf13/cobra"
)

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update the settings of an organization",
		Long:  `Update the settings of an organization`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req organizations.PatchedOrganization
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Organizations.UpdateOrganization(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Organization updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("org", "", "The organization name on Stackguardian platform.")
	_ = cmd.MarkFlagRequired("org")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
