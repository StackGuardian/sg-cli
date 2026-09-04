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

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create -- payload.json",
		Short: "Create an API access",
		Long:  `Create an API access`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req sggosdk.ApiAccess
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.AccessManagement.CreateApiAccess(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "API access created successfully.")
			return nil
		},
	}
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
