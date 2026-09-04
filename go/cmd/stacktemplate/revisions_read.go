// Code generated for sg-cli; edit the generator spec or this file directly.

package stacktemplate

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newRevisionsReadCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Get a stack template revision",
		Long:  `Get a stack template revision`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			revisionId := common.String(cmd, "revision-id")
			response, err := c.StackTemplateRevisions.ReadStackTemplateRevision(ctx, org, revisionId)
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().String("revision-id", "", "The revision ID (NAME:REVISION).")
	_ = cmd.MarkFlagRequired("revision-id")
	return cmd
}
