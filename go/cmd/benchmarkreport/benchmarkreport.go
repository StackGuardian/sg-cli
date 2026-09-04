// Code generated for sg-cli; edit the generator spec or this file directly.

package benchmarkreport

import (
	"github.com/StackGuardian/sg-cli/cmd/common"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

func NewBenchmarkReportCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmark-report",
		Short: "Read compliance benchmark reports",
		Long:  `Read compliance benchmark reports`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`Sub-commands:
  read                   Get the benchmark report`)
		},
	}
	common.AddOrgFlag(cmd)
	cmd.AddCommand(
		newReadCmd(c),
	)
	return cmd
}
