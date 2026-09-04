// Code generated for sg-cli; edit the generator spec or this file directly.

package extra

import (
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

// AddCommands attaches the generated sub-commands to an existing command group.
func AddCommands(cmd *cobra.Command, c *client.Client) {
	cmd.AddCommand(
		newReadCmd(c),
		newUpdateCmd(c),
		newListCmd(c),
		newCompareCmd(c),
	)
}
