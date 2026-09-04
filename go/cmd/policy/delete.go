// Code generated for sg-cli; edit the generator spec or this file directly.

package policy

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a policy",
		Long:  `Delete a policy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			policyId := common.String(cmd, "policy-id")
			if err := c.Policies.DeletePolicy(ctx, org, policyId); err != nil {
				return err
			}
			common.Done(cmd, nil, "Policy deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("policy-id", "", "The policy ID.")
	_ = cmd.MarkFlagRequired("policy-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
