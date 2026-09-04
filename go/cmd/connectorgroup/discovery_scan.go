// Code generated for sg-cli; edit the generator spec or this file directly.

package connectorgroup

import (
	"context"
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/connectorgroups"
	"github.com/spf13/cobra"
)

func newDiscoveryScanCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discovery-scan",
		Short: "Trigger discovery scans across the group's child connectors",
		Long:  `Trigger discovery scans across the group's child connectors`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			org := common.Org(cmd)
			connectorGroup := common.String(cmd, "connector-group")
			response, err := c.ConnectorGroups.TriggerDiscoveryScan(ctx, connectorGroup, org, &connectorgroups.DiscoveryScanRequest{Benchmark: common.OptString(cmd, "benchmark")})
			if err != nil {
				return err
			}
			common.Done(cmd, response, "Discovery scan triggered.")
			return nil
		},
	}
	cmd.Flags().String("connector-group", "", "The connector group name.")
	_ = cmd.MarkFlagRequired("connector-group")
	cmd.Flags().String("benchmark", "", "Benchmarks to scan (comma-separated, default inventory).")
	common.AddOutputJSONFlag(cmd)
	return cmd
}
