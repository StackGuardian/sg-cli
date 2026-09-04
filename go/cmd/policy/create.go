// Code generated for sg-cli; edit the generator spec or this file directly.

package policy

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
		Short: "Create a new policy (payload must include PolicyType: GENERAL or FILTER.INSIGHT)",
		Long:  `Create a new policy (payload must include PolicyType: GENERAL or FILTER.INSIGHT)`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			var req sggosdk.PolymorphicPolicy
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Policies.CreatePolicy(ctx, org, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Policy created successfully.")
			return nil
		},
	}
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
