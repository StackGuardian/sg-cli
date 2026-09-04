// Code generated for sg-cli; edit the generator spec or this file directly.

package apitoken

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Revoke the API key of a principal",
		Long:  `Revoke the API key of a principal`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			principalId := common.String(cmd, "principal-id")
			response, err := c.AccessManagement.DeleteApiToken(ctx, org, principalId)
			if err != nil {
				return err
			}
			common.Done(cmd, response, "API token revoked.")
			return nil
		},
	}
	cmd.Flags().String("principal-id", "", "Principal id in the form <userPoolId>/<idp>/<email>.")
	_ = cmd.MarkFlagRequired("principal-id")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
