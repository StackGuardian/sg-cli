package tests

import (
	"net/url"
	"testing"

	"github.com/StackGuardian/sg-cli/cmd/apiaccess"
	"github.com/StackGuardian/sg-cli/cmd/apitoken"
	"github.com/StackGuardian/sg-cli/cmd/auditlog"
	"github.com/StackGuardian/sg-cli/cmd/benchmarkreport"
	"github.com/StackGuardian/sg-cli/cmd/billing"
	"github.com/StackGuardian/sg-cli/cmd/chat"
	"github.com/StackGuardian/sg-cli/cmd/connector"
	"github.com/StackGuardian/sg-cli/cmd/connectorgroup"
	"github.com/StackGuardian/sg-cli/cmd/organization"
	"github.com/StackGuardian/sg-cli/cmd/policy"
	"github.com/StackGuardian/sg-cli/cmd/resource"
	"github.com/StackGuardian/sg-cli/cmd/role"
	"github.com/StackGuardian/sg-cli/cmd/rolebinding"
	"github.com/StackGuardian/sg-cli/cmd/runnergroup"
	"github.com/StackGuardian/sg-cli/cmd/secret"
	"github.com/StackGuardian/sg-cli/cmd/statebackend"
	"github.com/StackGuardian/sg-cli/cmd/user"
)

const org = "/api/v1/orgs/demo-org/"

func TestOrganizationCommands(t *testing.T) {
	runMatrix(t, []matrixCase{
		{name: "read", ctor: organization.NewOrganizationCmd, args: []string{"read", "--org", "demo-org"}, method: "GET", path: org},
		{name: "create", ctor: organization.NewOrganizationCmd, args: []string{"create", "@payload"}, payloadJSON: `{"ResourceName":"neworg"}`, method: "POST", path: "/api/v1/orgs/", body: `{"ResourceName":"neworg"}`, wantOut: "Organization created successfully."},
		{name: "list", ctor: organization.NewOrganizationCmd, args: []string{"list"}, method: "GET", path: "/api/v1/orgs/listall/", resp: `{"msg":["a"]}`},
		{name: "update", ctor: organization.NewOrganizationCmd, args: []string{"update", "--org", "demo-org", "@payload"}, payloadJSON: `{"ExtraSgCreditsLimit":5}`, method: "PATCH", path: org, body: `{"ExtraSgCreditsLimit":5}`, wantOut: "Organization updated successfully."},
		{name: "delete", ctor: organization.NewOrganizationCmd, args: []string{"delete", "--org", "demo-org"}, method: "DELETE", path: org, wantOut: "Organization archived successfully."},
		{name: "workflows", ctor: organization.NewOrganizationCmd, args: []string{"workflows", "--org", "demo-org", "--limit", "5", "--fetch-sg-owned-wfs=false", "--resource-names", "n"}, method: "GET", path: org + "wfs/listall/", query: url.Values{"limit": {"5"}, "fetchSgOwnedWfs": {"false"}, "resourceNames": {"n"}}, resp: `{"msg":[]}`},
		{name: "count", ctor: organization.NewOrganizationCmd, args: []string{"count-workflows", "--org", "demo-org", "--tags", "t"}, method: "GET", path: org + "wfs/listall/", query: url.Values{"onlyCount": {"1"}, "tags": {"t"}}, resp: `{"msg":3}`},
		{name: "bulk action", ctor: organization.NewOrganizationCmd, args: []string{"bulk-action", "--org", "demo-org", "--action", "activate", "--resource-ids", "/wfgrps/g/wfs/w,/wfgrps/g/stacks/s"}, method: "POST", path: org + "wfs/bulk-action/", body: `{"action":"activate","resourceIds":["/wfgrps/g/wfs/w","/wfgrps/g/stacks/s"]}`, wantOut: "Bulk action submitted."},
	})
}

func TestAccessCommands(t *testing.T) {
	runMatrix(t, []matrixCase{
		{name: "user read", ctor: user.NewUserCmd, args: []string{"read", "--org", "demo-org", "--user-id", "u@x"}, method: "POST", path: org + "get_user/", body: `{"userId":"u@x"}`},
		{name: "user create", ctor: user.NewUserCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"userId":"u@x","role":"ADMIN"}`, method: "POST", path: org + "invite_user/", body: `{"userId":"u@x","role":"ADMIN"}`, wantOut: "User invited successfully."},
		{name: "user delete", ctor: user.NewUserCmd, args: []string{"delete", "--org", "demo-org", "--alias", "grp"}, method: "POST", path: org + "remove_user/", body: `{"alias":"grp"}`, wantOut: "User removed successfully."},
		{name: "user update", ctor: user.NewUserCmd, args: []string{"update", "--org", "demo-org", "@payload"}, payloadJSON: `{"userId":"u@x","role":"USER"}`, method: "POST", path: org + "update_user/", body: `{"userId":"u@x","role":"USER"}`, wantOut: "User updated successfully."},
		{name: "user list", ctor: user.NewUserCmd, args: []string{"list", "--org", "demo-org", "--entity-type", "EMAIL", "--roles", "ADMIN"}, method: "GET", path: org + "users/listall/", query: url.Values{"entityType": {"EMAIL"}, "roles": {"ADMIN"}}},
		{name: "role create", ctor: role.NewRoleCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"r1"}`, method: "POST", path: org + "roles/", body: `{"ResourceName":"r1"}`, wantOut: "Role created successfully."},
		{name: "role read", ctor: role.NewRoleCmd, args: []string{"read", "--org", "demo-org", "--role-id", "r1"}, method: "GET", path: org + "roles/r1/"},
		{name: "role update", ctor: role.NewRoleCmd, args: []string{"update", "--org", "demo-org", "--role-id", "r1", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: org + "roles/r1/", body: `{"ResourceName":"r1","Description":"d"}`, wantOut: "Role updated successfully."},
		{name: "role delete", ctor: role.NewRoleCmd, args: []string{"delete", "--org", "demo-org", "--role-id", "r1"}, method: "DELETE", path: org + "roles/r1/", wantOut: "Role deleted successfully."},
		{name: "role list", ctor: role.NewRoleCmd, args: []string{"list", "--org", "demo-org", "--include-predefined-roles"}, method: "GET", path: org + "roles/listall/", query: url.Values{"includePredefinedRoles": {"true"}}, resp: `{"msg":[]}`},
		{name: "role binding", ctor: rolebinding.NewRoleBindingCmd, args: []string{"read", "--org", "demo-org"}, method: "GET", path: org + "rolebindings/default/"},
		{name: "api access create", ctor: apiaccess.NewApiAccessCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"k1"}`, method: "POST", path: org + "apiaccesses/", body: `{"ResourceName":"k1"}`, wantOut: "API access created successfully."},
		{name: "api access read", ctor: apiaccess.NewApiAccessCmd, args: []string{"read", "--org", "demo-org", "--access-id", "id1"}, method: "GET", path: org + "apiaccesses/id1/"},
		{name: "api access update", ctor: apiaccess.NewApiAccessCmd, args: []string{"update", "--org", "demo-org", "--access-id", "id1", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: org + "apiaccesses/id1/", body: `{"Description":"d"}`, wantOut: "API access updated successfully."},
		{name: "api access delete", ctor: apiaccess.NewApiAccessCmd, args: []string{"delete", "--org", "demo-org", "--access-id", "id1"}, method: "DELETE", path: org + "apiaccesses/id1/", wantOut: "API access deleted successfully."},
		{name: "api access regenerate", ctor: apiaccess.NewApiAccessCmd, args: []string{"regenerate-key", "--org", "demo-org", "--access-id", "id1", "--no-expiration"}, method: "POST", path: org + "apiaccesses/id1/regenerate/", body: `{"ExpiresAt":null}`, wantOut: "API key regenerated."},
		{name: "api access list", ctor: apiaccess.NewApiAccessCmd, args: []string{"list", "--org", "demo-org", "--access-type", "APIKEY", "--limit", "10", "--status", "active"}, method: "GET", path: org + "apiaccesses/listall/", query: url.Values{"accessType": {"APIKEY"}, "limit": {"10"}, "status": {"active"}}},
		{name: "api token create", ctor: apitoken.NewApiTokenCmd, args: []string{"create", "--org", "demo-org", "--regenerate"}, method: "POST", path: org + "api_token/", body: `{"regenerate":true}`, wantOut: "API token request completed."},
		{name: "api token delete", ctor: apitoken.NewApiTokenCmd, args: []string{"delete", "--org", "demo-org", "--principal-id", "pool/idp/u@x"}, method: "DELETE", path: org + "api_token/", body: `"pool/idp/u@x"`, wantOut: "API token revoked."},
		{name: "audit logs", ctor: auditlog.NewAuditLogCmd, args: []string{"list", "--org", "demo-org", "--effect", "Allow", "--start-time", "1", "--limit", "5"}, method: "GET", path: org + "audit_logs/", query: url.Values{"effect": {"Allow"}, "startTime": {"1"}, "limit": {"5"}}},
		{name: "benchmark", ctor: benchmarkreport.NewBenchmarkReportCmd, args: []string{"read", "--org", "demo-org", "--detailed", "--filter-csp", "AWS", "--limit", "2"}, method: "GET", path: org + "reports/benchmark/", query: url.Values{"detailed": {"true"}, "filter:CSP": {"AWS"}, "limit": {"2"}}, resp: `[]`},
	})
}

func TestSecretPolicyCommands(t *testing.T) {
	runMatrix(t, []matrixCase{
		{name: "secret create", ctor: secret.NewSecretCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"s","ResourceValue":"v"}`, method: "POST", path: org + "secrets/", body: `{"ResourceName":"s","ResourceValue":"v"}`, wantOut: "Secret created successfully."},
		{name: "secret update", ctor: secret.NewSecretCmd, args: []string{"update", "--org", "demo-org", "--secret-id", "s", "@payload"}, payloadJSON: `{"ResourceValue":"v2"}`, method: "PATCH", path: org + "secrets/s/", body: `{"ResourceName":"s","ResourceValue":"v2"}`, wantOut: "Secret updated successfully."},
		{name: "secret delete", ctor: secret.NewSecretCmd, args: []string{"delete", "--org", "demo-org", "--secret-id", "s"}, method: "DELETE", path: org + "secrets/s/", wantOut: "Secret deleted successfully."},
		{name: "secret list", ctor: secret.NewSecretCmd, args: []string{"list", "--org", "demo-org"}, method: "GET", path: org + "secrets/listall/"},
		{name: "secret read bulk", ctor: secret.NewSecretCmd, args: []string{"read-bulk", "--org", "demo-org", "--names", "a,b"}, method: "POST", path: org + "secrets/listall/", body: `{"ResourceNames":["a","b"]}`, resp: `{"msg":[]}`},
		{name: "policy create", ctor: policy.NewPolicyCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"PolicyType":"GENERAL","ResourceName":"p"}`, method: "POST", path: org + "policies/", body: `{"PolicyType":"GENERAL","ResourceName":"p"}`, wantOut: "Policy created successfully."},
		{name: "policy read", ctor: policy.NewPolicyCmd, args: []string{"read", "--org", "demo-org", "--policy-id", "p"}, method: "GET", path: org + "policies/p/"},
		{name: "policy update", ctor: policy.NewPolicyCmd, args: []string{"update", "--org", "demo-org", "--policy-id", "p", "@payload"}, payloadJSON: `{"PolicyType":"GENERAL","Description":"d"}`, method: "PATCH", path: org + "policies/p/", body: `{"PolicyType":"GENERAL","Description":"d"}`, wantOut: "Policy updated successfully."},
		{name: "policy delete", ctor: policy.NewPolicyCmd, args: []string{"delete", "--org", "demo-org", "--policy-id", "p"}, method: "DELETE", path: org + "policies/p/", wantOut: "Policy deleted successfully."},
		{name: "policy list", ctor: policy.NewPolicyCmd, args: []string{"list", "--org", "demo-org", "--search-query", "q", "--limit", "5"}, method: "GET", path: org + "policies/listall/", query: url.Values{"searchQuery": {"q"}, "limit": {"5"}}, resp: `{"msg":[]}`},
	})
}

func TestConnectorCommands(t *testing.T) {
	base := org + "integrations/"
	gbase := org + "integrationgroups/"
	runMatrix(t, []matrixCase{
		{name: "create", ctor: connector.NewConnectorCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"c"}`, method: "POST", path: base, body: `{"ResourceName":"c"}`, wantOut: "Connector created successfully."},
		{name: "read", ctor: connector.NewConnectorCmd, args: []string{"read", "--org", "demo-org", "--connector-id", "c"}, method: "GET", path: base + "c/"},
		{name: "update", ctor: connector.NewConnectorCmd, args: []string{"update", "--org", "demo-org", "--connector-id", "c", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: base + "c/", body: `{"Description":"d"}`, wantOut: "Connector updated successfully."},
		{name: "delete", ctor: connector.NewConnectorCmd, args: []string{"delete", "--org", "demo-org", "--connector-id", "c"}, method: "DELETE", path: base + "c/", wantOut: "Connector deleted successfully."},
		{name: "list", ctor: connector.NewConnectorCmd, args: []string{"list", "--org", "demo-org", "--kind", "AWS_RBAC", "--limit", "5"}, method: "GET", path: base + "listall/", query: url.Values{"kind": {"AWS_RBAC"}, "limit": {"5"}}},
		{name: "authenticate", ctor: connector.NewConnectorCmd, args: []string{"authenticate", "--org", "demo-org", "--connector-id", "c"}, method: "GET", path: base + "c/authenticate/"},
		{name: "list accounts", ctor: connector.NewConnectorCmd, args: []string{"list-accounts", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceId":"g"}`, method: "POST", path: base + "list_accounts/", body: `{"ResourceId":"g"}`, resp: `{"msg":[]}`, wantOut: `"msg": []`},
		{name: "github repos", ctor: connector.NewConnectorCmd, args: []string{"github-repos", "--org", "demo-org", "--connector-id", "c"}, method: "GET", path: base + "c/get_githubcom_repos/", resp: `{"msg":["o/r"]}`},
		{name: "repos", ctor: connector.NewConnectorCmd, args: []string{"repos", "--org", "demo-org", "--connector-id", "c", "--login", "l", "--page", "2", "--limit", "5"}, method: "GET", path: base + "c/repos/", query: url.Values{"login": {"l"}, "page": {"2"}, "limit": {"5"}}},
		{name: "group create", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"g"}`, method: "POST", path: gbase, body: `{"ResourceName":"g"}`, wantOut: "Connector group created successfully."},
		{name: "group read", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"read", "--org", "demo-org", "--connector-group", "g"}, method: "GET", path: gbase + "g/"},
		{name: "group update", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"update", "--org", "demo-org", "--connector-group", "g", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: gbase + "g/", body: `{"Description":"d"}`, wantOut: "Connector group updated successfully."},
		{name: "group delete", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"delete", "--org", "demo-org", "--connector-group", "g"}, method: "DELETE", path: gbase + "g/", wantOut: "Connector group deleted successfully."},
		{name: "group authenticate", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"authenticate", "--org", "demo-org", "--connector-group", "g"}, method: "GET", path: gbase + "g/authenticate/"},
		{name: "group list", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"list", "--org", "demo-org", "--scope", "s"}, method: "GET", path: gbase + "listall/", query: url.Values{"scope": {"s"}}},
		{name: "group discovery scan", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"discovery-scan", "--org", "demo-org", "--connector-group", "g", "--benchmark", "inventory"}, method: "POST", path: gbase + "g/discovery_scan/", query: url.Values{"benchmark": {"inventory"}}, wantOut: "Discovery scan triggered."},
		{name: "child read", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"connectors", "read", "--org", "demo-org", "--connector-group", "g", "--connector-id", "c"}, method: "GET", path: gbase + "g/integrations/c/"},
		{name: "child update", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"connectors", "update", "--org", "demo-org", "--connector-group", "g", "--connector-id", "c", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: gbase + "g/integrations/c/", body: `{"Description":"d"}`, wantOut: "Connector updated successfully."},
		{name: "child delete", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"connectors", "delete", "--org", "demo-org", "--connector-group", "g", "--connector-id", "c"}, method: "DELETE", path: gbase + "g/integrations/c/", wantOut: "Connector deleted successfully."},
		{name: "child list", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"connectors", "list", "--org", "demo-org", "--connector-group", "g", "--kind", "k"}, method: "GET", path: gbase + "g/integrations/listall/", query: url.Values{"kind": {"k"}}},
		{name: "child authenticate", ctor: connectorgroup.NewConnectorGroupCmd, args: []string{"connectors", "authenticate", "--org", "demo-org", "--connector-group", "g", "--connector-id", "c"}, method: "GET", path: gbase + "g/integrations/c/authenticate/"},
	})
}

func TestRunnerStateBackendResourceCommands(t *testing.T) {
	rbase := org + "runnergroups/"
	sbase := org + "statebackends/"
	runMatrix(t, []matrixCase{
		{name: "runner create", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"rg"}`, method: "POST", path: rbase, body: `{"ResourceName":"rg"}`, wantOut: "Runner group created successfully."},
		{name: "runner read", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"read", "--org", "demo-org", "--runner-group", "rg", "--get-active-workflows"}, method: "GET", path: rbase + "rg/", query: url.Values{"getActiveWorkflows": {"true"}}},
		{name: "runner update", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"update", "--org", "demo-org", "--runner-group", "rg", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: rbase + "rg/", body: `{"Description":"d"}`, wantOut: "Runner group updated successfully."},
		{name: "runner delete", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"delete", "--org", "demo-org", "--runner-group", "rg"}, method: "DELETE", path: rbase + "rg/", wantOut: "Runner group deleted successfully."},
		{name: "runner list", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"list", "--org", "demo-org", "--limit", "5"}, method: "GET", path: rbase + "listall/", query: url.Values{"limit": {"5"}}, resp: `{"msg":[]}`},
		{name: "runner register", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"register", "--org", "demo-org", "--runner-group", "rg"}, method: "POST", path: rbase + "rg/register/", wantOut: "Runner registered."},
		{name: "runner storage auth", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"storage-backend-auth", "--org", "demo-org", "--runner-group", "rg"}, method: "POST", path: rbase + "rg/storage-backend-auth/"},
		{name: "runner deregister", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"deregister-runner", "--org", "demo-org", "--runner-group", "rg", "--runner-id", "r1", "--container-instance-ids", "a,b", "--force"}, method: "POST", path: rbase + "rg/deregister/", body: `{"RunnerId":"r1","ContainerInstanceIds":["a","b"],"ForceDeregister":true}`, wantOut: "Runner deregistered."},
		{name: "runner state", ctor: runnergroup.NewRunnerGroupCmd, args: []string{"update-runner-state", "--org", "demo-org", "--runner-group", "rg", "--status", "DRAINING", "--runner-id", "r1"}, method: "POST", path: rbase + "rg/runner_status/", body: `{"RunnerId":"r1","Status":"DRAINING"}`, wantOut: "Runner state updated."},
		{name: "state backend create", ctor: statebackend.NewStateBackendCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"sb","Statefiles":[],"StateBackendConfig":{"type":"aws_s3","s3BucketName":"b","auth":{"integrationId":"/integrations/aws"}}}`, method: "POST", path: sbase, body: `{"ResourceName":"sb","Statefiles":[],"StateBackendConfig":{"type":"aws_s3","s3BucketName":"b","auth":{"integrationId":"/integrations/aws"}}}`, wantOut: "State backend created successfully."},
		{name: "state backend read", ctor: statebackend.NewStateBackendCmd, args: []string{"read", "--org", "demo-org", "--state-backend", "sb"}, method: "GET", path: sbase + "sb/"},
		{name: "state backend update", ctor: statebackend.NewStateBackendCmd, args: []string{"update", "--org", "demo-org", "--state-backend", "sb", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: sbase + "sb/", body: `{"Description":"d"}`, wantOut: "State backend updated successfully."},
		{name: "state backend delete", ctor: statebackend.NewStateBackendCmd, args: []string{"delete", "--org", "demo-org", "--state-backend", "sb"}, method: "DELETE", path: sbase + "sb/", wantOut: "State backend deleted successfully."},
		{name: "state backend list", ctor: statebackend.NewStateBackendCmd, args: []string{"list", "--org", "demo-org", "--search-query", "q"}, method: "GET", path: sbase + "listall/", query: url.Values{"searchQuery": {"q"}}},
		{name: "state backend statefiles", ctor: statebackend.NewStateBackendCmd, args: []string{"list-statefiles", "--org", "demo-org", "--config", "abc"}, method: "GET", path: sbase + "list_statefiles/", query: url.Values{"config": {"abc"}}},
		{name: "resource move", ctor: resource.NewResourceCmd, args: []string{"move", "--org", "demo-org", "--source-ids", "/a,/b", "--destination", "/wfgrps/g"}, method: "POST", path: org + "resources/move/", body: `{"SourceResourceIds":["/a","/b"],"DestinationParentId":"/wfgrps/g"}`, wantOut: "Resources moved."},
		{name: "resource search", ctor: resource.NewResourceCmd, args: []string{"search", "--org", "demo-org", "--search-query", "q", "--resource-types", "WORKFLOW,STACK", "--limit", "5"}, method: "POST", path: org + "search/", body: `{"SearchQuery":"q","ResourceTypes":["WORKFLOW","STACK"],"Limit":5}`, wantOut: `"msg"`},
		{name: "resource tags", ctor: resource.NewResourceCmd, args: []string{"tags", "--org", "demo-org", "--context-tag-key", "env", "--limit", "5"}, method: "GET", path: org + "search/tags/", query: url.Values{"ContextTagKey": {"env"}, "Limit": {"5"}}},
	})
}

func TestBillingChatCommands(t *testing.T) {
	bbase := org + "billing/"
	cbase := org + "chats/"
	runMatrix(t, []matrixCase{
		{name: "balance", ctor: billing.NewBillingCmd, args: []string{"balance", "--org", "demo-org"}, method: "GET", path: bbase + "list-balance/"},
		{name: "dashboard", ctor: billing.NewBillingCmd, args: []string{"dashboard-url", "--org", "demo-org", "--type", "usage"}, method: "GET", path: bbase + "dashboard/", query: url.Values{"type": {"usage"}}},
		{name: "invoices", ctor: billing.NewBillingCmd, args: []string{"invoices", "--org", "demo-org", "--limit", "5", "--next-page", "c"}, method: "GET", path: bbase + "invoices/", query: url.Values{"limit": {"5"}, "next_page": {"c"}}},
		{name: "details", ctor: billing.NewBillingCmd, args: []string{"details", "--org", "demo-org"}, method: "GET", path: bbase + "details/"},
		{name: "change plan", ctor: billing.NewBillingCmd, args: []string{"change-plan", "--org", "demo-org", "--target-package", "SG | Standard"}, method: "POST", path: bbase + "change-plan/", body: `{"TargetPackage":"SG | Standard"}`, wantOut: "Plan change submitted."},
		{name: "profile setup", ctor: billing.NewBillingCmd, args: []string{"profile", "setup", "--org", "demo-org", "@payload"}, payloadJSON: `{"Email":"e@x"}`, method: "POST", path: bbase + "profile/", body: `{"Email":"e@x"}`, wantOut: "Billing profile created."},
		{name: "profile update", ctor: billing.NewBillingCmd, args: []string{"profile", "update", "--org", "demo-org", "@payload"}, payloadJSON: `{"Name":"n"}`, method: "PATCH", path: bbase + "profile/", body: `{"Name":"n"}`, wantOut: "Billing profile updated."},
		{name: "payment list", ctor: billing.NewBillingCmd, args: []string{"payment-methods", "list", "--org", "demo-org"}, method: "GET", path: bbase + "payment-methods/"},
		{name: "payment intent", ctor: billing.NewBillingCmd, args: []string{"payment-methods", "setup-intent", "--org", "demo-org"}, method: "POST", path: bbase + "payment-methods/"},
		{name: "payment default", ctor: billing.NewBillingCmd, args: []string{"payment-methods", "set-default", "--org", "demo-org", "--payment-method-id", "pm_1"}, method: "POST", path: bbase + "set-default-payment-method/", body: `{"PaymentMethodId":"pm_1"}`, wantOut: "Default payment method set."},
		{name: "payment detach", ctor: billing.NewBillingCmd, args: []string{"payment-methods", "detach", "--org", "demo-org", "--payment-method-id", "pm_1"}, method: "POST", path: bbase + "detach-card/", body: `{"PaymentMethodId":"pm_1"}`, wantOut: "Card detached."},
		{name: "chat create", ctor: chat.NewChatCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ChatType":"kyro","Config":{}}`, method: "POST", path: cbase, body: `{"ChatType":"kyro","Config":{}}`, wantOut: "Chat created successfully."},
		{name: "chat list", ctor: chat.NewChatCmd, args: []string{"list", "--org", "demo-org", "--chat-type", "kyro", "--refresh-pr-status"}, method: "GET", path: cbase + "listall/", query: url.Values{"chatType": {"kyro"}, "refreshPrStatus": {"true"}}, resp: `{"msg":[]}`},
		{name: "chat read", ctor: chat.NewChatCmd, args: []string{"read", "--org", "demo-org", "--chat-id", "c1"}, method: "GET", path: cbase + "c1/"},
		{name: "chat update", ctor: chat.NewChatCmd, args: []string{"update", "--org", "demo-org", "--chat-id", "c1", "@payload"}, payloadJSON: `{"Status":"archived"}`, method: "PATCH", path: cbase + "c1/", body: `{"Status":"archived"}`, wantOut: "Chat updated successfully."},
		{name: "chat delete", ctor: chat.NewChatCmd, args: []string{"delete", "--org", "demo-org", "--chat-id", "c1"}, method: "DELETE", path: cbase + "c1/", wantOut: "Chat archived successfully."},
		{name: "chat commit", ctor: chat.NewChatCmd, args: []string{"commit-bundle", "--org", "demo-org", "--chat-id", "c1", "--expected-version", "3"}, method: "POST", path: cbase + "c1/commit_bundle/", body: `{"expected_version":3}`, wantOut: "Bundle committed."},
		{name: "chat create pr", ctor: chat.NewChatCmd, args: []string{"create-pr", "--org", "demo-org", "--chat-id", "c1", "@payload"}, payloadJSON: `{"vcsConnector":"v","repoUrl":"https://x"}`, method: "POST", path: cbase + "c1/create_pr/", body: `{"vcsConnector":"v","repoUrl":"https://x"}`, wantOut: "Pull request request completed."},
		{name: "chat pr status", ctor: chat.NewChatCmd, args: []string{"pr-status", "--org", "demo-org", "--chat-id", "c1"}, method: "GET", path: cbase + "c1/pr_status/"},
		{name: "chat pr sync", ctor: chat.NewChatCmd, args: []string{"pr-sync-status", "--org", "demo-org", "--chat-id", "c1"}, method: "GET", path: cbase + "c1/pr_sync_status/"},
		{name: "chat artifacts list", ctor: chat.NewChatCmd, args: []string{"artifacts", "list", "--org", "demo-org", "--chat-id", "c1", "--prefix", "p"}, method: "GET", path: cbase + "c1/listall_artifacts/", query: url.Values{"prefix": {"p"}}},
		{name: "chat artifact url", ctor: chat.NewChatCmd, args: []string{"artifacts", "get-url", "--org", "demo-org", "--chat-id", "c1", "--artifact-path", "a"}, method: "GET", path: cbase + "c1/get_artifact/", query: url.Values{"artifactPath": {"a"}}},
		{name: "chat artifact upload", ctor: chat.NewChatCmd, args: []string{"artifacts", "upload-url", "--org", "demo-org", "--chat-id", "c1", "--artifact-path", "a", "--content-type", "text/plain"}, method: "POST", path: cbase + "c1/save_artifact/", query: url.Values{"artifactPath": {"a"}, "contentType": {"text/plain"}}},
		{name: "message create", ctor: chat.NewChatCmd, args: []string{"messages", "create", "--org", "demo-org", "--chat-id", "c1", "--text", "hi", "--client-request-id", "r1"}, method: "POST", path: cbase + "c1/messages/", body: `{"content":[{"type":"text","text":"hi"}],"clientRequestId":"r1"}`, wantOut: "Message sent."},
		{name: "message create payload", ctor: chat.NewChatCmd, args: []string{"messages", "create", "--org", "demo-org", "--chat-id", "c1", "@payload"}, payloadJSON: `{"content":[{"type":"tool_result","toolUseId":"t","isError":true}]}`, method: "POST", path: cbase + "c1/messages/", body: `{"content":[{"type":"tool_result","toolUseId":"t","isError":true}]}`, wantOut: "Message sent."},
		{name: "messages list", ctor: chat.NewChatCmd, args: []string{"messages", "list", "--org", "demo-org", "--chat-id", "c1", "--limit", "5"}, method: "GET", path: cbase + "c1/messages/listall/", query: url.Values{"limit": {"5"}}, resp: `{"msg":[]}`},
		{name: "message read", ctor: chat.NewChatCmd, args: []string{"messages", "read", "--org", "demo-org", "--chat-id", "c1", "--message-id", "m1"}, method: "GET", path: cbase + "c1/messages/m1/"},
		{name: "message cancel", ctor: chat.NewChatCmd, args: []string{"messages", "cancel", "--org", "demo-org", "--chat-id", "c1", "--message-id", "m1"}, method: "POST", path: cbase + "c1/messages/m1/cancel/", wantOut: "Message cancelled."},
		{name: "message retry", ctor: chat.NewChatCmd, args: []string{"messages", "retry", "--org", "demo-org", "--chat-id", "c1", "--message-id", "m1"}, method: "POST", path: cbase + "c1/messages/m1/retry/", wantOut: "Message retried."},
	})
}
