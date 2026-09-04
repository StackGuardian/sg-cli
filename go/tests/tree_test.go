package tests

import (
	"sort"
	"strings"
	"testing"

	"github.com/StackGuardian/sg-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// expectedLeaves pins the complete command surface (one leaf per SDK method plus
// the apply/destroy conveniences). Update it deliberately when commands change.
const expectedLeaves = `
api-access create
api-access delete
api-access list
api-access read
api-access regenerate-key
api-access update
api-token create
api-token delete
artifacts create
artifacts delete
artifacts get-url
artifacts list
artifacts lock
artifacts read
artifacts rollback
artifacts unlock
artifacts versions
audit-log list
benchmark-report read
billing balance
billing change-plan
billing dashboard-url
billing details
billing invoices
billing payment-methods detach
billing payment-methods list
billing payment-methods set-default
billing payment-methods setup-intent
billing profile setup
billing profile update
chat artifacts get-url
chat artifacts list
chat artifacts upload-url
chat commit-bundle
chat create
chat create-pr
chat delete
chat list
chat messages cancel
chat messages create
chat messages list
chat messages read
chat messages retry
chat pr-status
chat pr-sync-status
chat read
chat update
connector authenticate
connector create
connector delete
connector github-repos
connector list
connector list-accounts
connector read
connector repos
connector update
connector-group authenticate
connector-group connectors authenticate
connector-group connectors delete
connector-group connectors list
connector-group connectors read
connector-group connectors update
connector-group create
connector-group delete
connector-group discovery-scan
connector-group list
connector-group read
connector-group update
organization bulk-action
organization count-workflows
organization create
organization delete
organization list
organization read
organization update
organization workflows
policy create
policy delete
policy list
policy read
policy update
resource move
resource search
resource tags
role create
role delete
role list
role read
role update
role-binding read
runner-group create
runner-group delete
runner-group deregister-runner
runner-group list
runner-group read
runner-group register
runner-group storage-backend-auth
runner-group update
runner-group update-runner-state
secret create
secret delete
secret list
secret read-bulk
secret update
stack apply
stack compare
stack create
stack delete
stack destroy
stack list
stack outputs
stack read
stack update
stack-run create
stack-run list
stack-run read
stack-template create
stack-template delete
stack-template read
stack-template revisions create
stack-template revisions delete
stack-template revisions read
stack-template revisions update
stack-template update
stack-workflow artifacts create
stack-workflow artifacts delete
stack-workflow artifacts get-url
stack-workflow artifacts list
stack-workflow artifacts lock
stack-workflow artifacts read
stack-workflow artifacts unlock
stack-workflow compare
stack-workflow create
stack-workflow delete
stack-workflow file-upload-url
stack-workflow list
stack-workflow outputs
stack-workflow read
stack-workflow tfstate-upload-url
stack-workflow update
stack-workflow vcs-triggers
stack-workflow-run approve
stack-workflow-run create
stack-workflow-run delete
stack-workflow-run facts create
stack-workflow-run facts read
stack-workflow-run facts update
stack-workflow-run list
stack-workflow-run logs
stack-workflow-run read
stack-workflow-run update
state-backend create
state-backend delete
state-backend list
state-backend list-statefiles
state-backend read
state-backend update
template artifacts delete
template artifacts download-url
template artifacts list
template artifacts upload-url
template create
template delete
template group-iac
template input-schema
template list
template list-by-owner
template public list
template public list-all
template public read
template read
template subscriptions create
template subscriptions read
template subscriptions update
template update
template vcs-triggers
user create
user delete
user list
user read
user update
workflow apply
workflow compare
workflow create
workflow delete
workflow destroy
workflow file-upload-url
workflow list
workflow outputs
workflow read
workflow tfstate-upload-url
workflow update
workflow vcs-triggers
workflow-group create
workflow-group create-child
workflow-group delete
workflow-group list
workflow-group list-children
workflow-group list-resources
workflow-group read
workflow-group update
workflow-run approve
workflow-run cancel
workflow-run create
workflow-run delete
workflow-run facts create
workflow-run facts read
workflow-run facts update
workflow-run list
workflow-run logs
workflow-run read
workflow-run read-by-ksuid
workflow-run update
workflow-run update-by-ksuid
workflow-step-template create
workflow-step-template delete
workflow-step-template read
workflow-step-template revisions create
workflow-step-template revisions delete
workflow-step-template revisions read
workflow-step-template revisions update
workflow-step-template update
workflow-template create
workflow-template delete
workflow-template read
workflow-template revisions create
workflow-template revisions delete
workflow-template revisions read
workflow-template revisions update
workflow-template update
`

func leafCommands(root *cobra.Command) []string {
	var leaves []string
	var walk func(c *cobra.Command, prefix string)
	walk = func(c *cobra.Command, prefix string) {
		if c.Name() == "help" || c.Name() == "completion" {
			return
		}
		if !c.HasSubCommands() {
			leaves = append(leaves, prefix+c.Name())
			return
		}
		for _, sub := range c.Commands() {
			walk(sub, prefix+c.Name()+" ")
		}
	}
	for _, sub := range root.Commands() {
		walk(sub, "")
	}
	sort.Strings(leaves)
	return leaves
}

func TestCommandTree(t *testing.T) {
	want := strings.Split(strings.TrimSpace(expectedLeaves), "\n")
	assert.Equal(t, want, leafCommands(cmd.RootCommand()))
	assert.Len(t, want, 230)
}
