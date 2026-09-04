// Code generated for sg-cli; edit the generator spec or this file directly.

package stacktemplate

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newRevisionsDeleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a stack template revision",
		Long:  `Delete a stack template revision`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			revisionId := common.String(cmd, "revision-id")
			keepParentTemplate := common.Bool(cmd, "keep-parent-template")
			if err := c.StackTemplateRevisions.DeleteStackTemplateRevision(ctx, org, revisionId, keepParentTemplate); err != nil {
				return err
			}
			common.Done(cmd, nil, "Revision deleted successfully.")
			return nil
		},
	}
	cmd.Flags().String("revision-id", "", "The revision ID (NAME:REVISION).")
	_ = cmd.MarkFlagRequired("revision-id")
	cmd.Flags().Bool("keep-parent-template", false, "Keep the parent template when deleting its last revision.")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
