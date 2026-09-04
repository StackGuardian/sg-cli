// Code generated for sg-cli; edit the generator spec or this file directly.

package auditlog

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewAuditLogCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit-log",
		Short: "Read organization audit logs",
		Long:  `Read organization audit logs`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  list                   List audit log entries`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newListCmd(c),
	)
	return cmd
}
