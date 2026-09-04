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

func newUpdateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update -- payload.json",
		Short: "Update a policy (payload must include PolicyType)",
		Long:  `Update a policy (payload must include PolicyType)`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			policyId := common.String(cmd, "policy-id")
			var req sggosdk.PatchedPolymorphicPolicy
			if err := common.ReadPayload(cmd, args, &req); err != nil {
				if errors.Is(err, common.ErrDryRun) {
					return nil
				}
				return err
			}
			response, err := c.Policies.UpdatePolicy(ctx, org, policyId, &req)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Policy updated successfully.")
			return nil
		},
	}
	cmd.Flags().String("policy-id", "", "The policy ID.")
	_ = cmd.MarkFlagRequired("policy-id")
	common.AddPayloadFlags(cmd)
	common.AddOutputJSONFlag(cmd)
	return cmd
}
