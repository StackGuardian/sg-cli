// Code generated for sg-cli; edit the generator spec or this file directly.

package user

import (
	"context"
	"errors"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Invite a user or SSO group to the organization",
		Long:  `Invite a user or SSO group to the organization`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req sggosdk.AddUserToOrganization
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.AccessManagement.CreateUser(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "User invited successfully.")
			return nil
		},
	}
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
