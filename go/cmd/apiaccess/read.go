// Code generated for sg-cli; edit the generator spec or this file directly.

package apiaccess

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get an API access",
		Long:  `Get an API access`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			accessId := common.String(cmd, "access-id")
			response, err := c.AccessManagement.ReadApiAccess(ctx, accessId, org)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("access-id", "", "The API access ID.")
	_ = cmd.MarkFlagRequired("access-id")
	return cmd
}
