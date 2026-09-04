// Code generated for sg-cli; edit the generator spec or this file directly.

package role

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a role",
		Long:  `Update a role`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			roleId := common.String(cmd, "role-id")
			var req sggosdk.PatchedRole
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			if req.ResourceName == nil {
				req.ResourceName = sggosdk.Optional(roleId) // the API requires it on updates
			}
			response, err := c.AccessManagement.UpdateRole(ctx, org, roleId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Role updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("role-id", "", "The role ID.")
	_ = cmd.MarkFlagRequired("role-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
