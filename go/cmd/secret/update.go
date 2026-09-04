// Code generated for sg-cli; edit the generator spec or this file directly.

package secret

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
		Short: "Update a secret",
		Long:  `Update a secret`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			secretId := common.String(cmd, "secret-id")
			var req sggosdk.PatchedSecret
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			if req.ResourceName == nil {
				req.ResourceName = sggosdk.Optional(secretId) // the API requires it on updates
			}
			response, err := c.Secrets.UpdateSecret(ctx, org, secretId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Secret updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("secret-id", "", "The secret name.")
	_ = cmd.MarkFlagRequired("secret-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
