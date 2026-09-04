// Code generated for sg-cli; edit the generator spec or this file directly.

package apiaccess

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
		Short: "Update an API access",
		Long:  `Update an API access`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			accessId := common.String(cmd, "access-id")
			var req sggosdk.PatchedApiAccessPatch
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.AccessManagement.UpdateApiAccess(ctx, accessId, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "API access updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("access-id", "", "The API access ID.")
	_ = cmd.MarkFlagRequired("access-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
