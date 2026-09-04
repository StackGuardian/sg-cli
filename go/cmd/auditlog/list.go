// Code generated for sg-cli; edit the generator spec or this file directly.

package auditlog

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func newListCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit log entries",
		Long:  `List audit log entries`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			response, err := c.AccessManagement.ReadAuditLogs(ctx, org, &sggosdk.ReadAuditLogsRequest{Effect: common.OptString(cmd, "effect"), EndTime: common.OptInt(cmd, "end-time"), Lastevaluatedkey: common.OptString(cmd, "last-evaluated-key"), Limit: common.OptInt(cmd, "limit"), PrincipalEmail: common.OptString(cmd, "principal-email"), RequestMethod: common.OptString(cmd, "request-method"), Resource: common.OptString(cmd, "resource"), SourceIp: common.OptString(cmd, "source-ip"), StartTime: common.OptInt(cmd, "start-time")})
			if err != nil {
				return err
			}
			common.PrintJSON(cmd, response)
			return nil
		},
	}
	cmd.Flags().Int("limit", 0, "Limit the number of results returned.")
	cmd.Flags().String("last-evaluated-key", "", "Pagination token from a previous response.")
	cmd.Flags().String("effect", "", "Filter by effect (Allow, Deny).")
	cmd.Flags().Int("start-time", 0, "Start time as a Unix timestamp in milliseconds.")
	cmd.Flags().Int("end-time", 0, "End time as a Unix timestamp in milliseconds.")
	cmd.Flags().String("principal-email", "", "Filter by principal email.")
	cmd.Flags().String("request-method", "", "Filter by HTTP method.")
	cmd.Flags().String("resource", "", "Filter by resource name.")
	cmd.Flags().String("source-ip", "", "Filter by source IP address.")
	return cmd
}
