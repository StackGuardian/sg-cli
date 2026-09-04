// Code generated for sg-cli; edit the generator spec or this file directly.

package apitoken

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/accessmanagement"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newCreateCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create or retrieve the caller's API key, or look up a runner group's key",
		Long:  `Create or retrieve the caller's API key, or look up a runner group's key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.AccessManagement.CreateApiToken(ctx, org, &accessmanagement.ApiTokenRequest{Regenerate: common.OptBool(cmd, "regenerate"), RunnerGroupId: common.OptString(cmd, "runner-group-id")})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "API token request completed.")
			return nil
		},
	}
	cmd.Flags().Bool("regenerate", false, "Replace the existing key with a new one.")
	cmd.Flags().String("runner-group-id", "", "Look up the key of this runner group (/runnergroups/<name>).")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
