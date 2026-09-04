package tests

import (
	"net/url"
	"testing"

	artifactscmd "github.com/StackGuardian/sg-cli/cmd/artifacts"
	stackcmd "github.com/StackGuardian/sg-cli/cmd/stack"
	"github.com/StackGuardian/sg-cli/cmd/stackrun"
	"github.com/StackGuardian/sg-cli/cmd/stacktemplate"
	"github.com/StackGuardian/sg-cli/cmd/stackworkflow"
	"github.com/StackGuardian/sg-cli/cmd/stackworkflowrun"
	"github.com/StackGuardian/sg-cli/cmd/template"
	workflowcmd "github.com/StackGuardian/sg-cli/cmd/workflow"
	"github.com/StackGuardian/sg-cli/cmd/workflowgroup"
	"github.com/StackGuardian/sg-cli/cmd/workflowrun"
	"github.com/StackGuardian/sg-cli/cmd/workflowsteptemplate"
	"github.com/StackGuardian/sg-cli/cmd/workflowtemplate"
)

func TestWorkflowGroupCommands(t *testing.T) {
	base := org + "wfgrps/"
	runMatrix(t, []matrixCase{
		{name: "create", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"ResourceName":"g"}`, method: "POST", path: base, body: `{"ResourceName":"g"}`, wantOut: "Workflow group created successfully."},
		{name: "read nested", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"read", "--org", "demo-org", "--workflow-group", "parent/child"}, method: "GET", path: base + "parent/child/"},
		{name: "update", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"update", "--org", "demo-org", "--workflow-group", "g", "@payload"}, payloadJSON: `{"Description":"d"}`, method: "PATCH", path: base + "g/", body: `{"Description":"d"}`, wantOut: "Workflow group updated successfully."},
		{name: "delete", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"delete", "--org", "demo-org", "--workflow-group", "g"}, method: "DELETE", path: base + "g/", wantOut: "Workflow group deleted successfully."},
		{name: "create child", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"create-child", "--org", "demo-org", "--parent", "g", "@payload"}, payloadJSON: `{"ResourceName":"c"}`, method: "POST", path: base + "g/wfgrps/", body: `{"ResourceName":"c"}`, wantOut: "Child workflow group created successfully."},
		{name: "list children", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"list-children", "--org", "demo-org", "--parent", "g", "--limit", "5"}, method: "GET", path: base + "g/wfgrps/listall/", query: url.Values{"limit": {"5"}}},
		{name: "list", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"list", "--org", "demo-org", "--search-query", "q"}, method: "GET", path: base + "listall/", query: url.Values{"searchQuery": {"q"}}},
		{name: "list resources", ctor: workflowgroup.NewWorkflowGroupCmd, args: []string{"list-resources", "--org", "demo-org", "--workflow-group", "g", "--resource-types", "WORKFLOW,STACK", "--limit", "5"}, method: "GET", path: base + "g/listall/", query: url.Values{"resourceTypes": {"WORKFLOW,STACK"}, "limit": {"5"}}, resp: `{"msg":[]}`},
	})
}

func TestWorkflowExtraCommands(t *testing.T) {
	base := org + "wfgrps/g/wfs/w/"
	scope := []string{"--org", "demo-org", "--workflow-group", "g"}
	with := func(args ...string) []string { return append(args, scope...) }
	runMatrix(t, []matrixCase{
		{name: "update", ctor: workflowcmd.NewWorkflowCmd, args: with("update", "--workflow-id", "w", "--upgrade-mode", "MANUAL", "@payload"), payloadJSON: `{"Description":"d"}`, method: "PATCH", path: base, query: url.Values{"upgradeMode": {"MANUAL"}}, body: `{"Description":"d"}`, wantOut: "Workflow updated successfully."},
		{name: "outputs", ctor: workflowcmd.NewWorkflowCmd, args: with("outputs", "--workflow-id", "w"), method: "GET", path: base + "outputs/"},
		{name: "compare", ctor: workflowcmd.NewWorkflowCmd, args: with("compare", "--workflow-id", "w", "--target-template-id", "/o/t"), method: "POST", path: base + "compare/", body: `{"targetTemplateId":"/o/t"}`, wantOut: `"msg"`},
		{name: "file upload url", ctor: workflowcmd.NewWorkflowCmd, args: with("file-upload-url", "--workflow-id", "w", "--filename", "f", "--folder", "d"), method: "GET", path: base + "file_upload_url/", query: url.Values{"filename": {"f"}, "folder": {"d"}}},
		{name: "tfstate upload url", ctor: workflowcmd.NewWorkflowCmd, args: with("tfstate-upload-url", "--workflow-id", "w", "--filename", "s.json"), method: "GET", path: base + "tfstate_upload_url/", query: url.Values{"filename": {"s.json"}}},
		{name: "vcs triggers", ctor: workflowcmd.NewWorkflowCmd, args: with("vcs-triggers", "--workflow-id", "w", "@payload"), payloadJSON: `{}`, method: "POST", path: base + "webhooks/vcs_triggers/", body: `{"VCSConfig":null,"VCSTriggers":null}`, wantOut: "VCS triggers created successfully."},
		{name: "artifact url", ctor: artifactscmd.NewArtifactsCmd, args: with("get-url", "--workflow-id", "w", "--artifact-path", "p", "--version-id", "v"), method: "GET", path: base + "get_artifact/", query: url.Values{"artifactPath": {"p"}, "versionId": {"v"}}},
		{name: "artifact read", ctor: artifactscmd.NewArtifactsCmd, args: with("read", "--workflow-id", "w", "--artifact-id", "a/b.json"), method: "GET", path: base + "artifacts/a/b.json/", resp: `{"v":4}`},
		{name: "artifact create", ctor: artifactscmd.NewArtifactsCmd, args: with("create", "--workflow-id", "w", "--artifact-id", "tf", "@payload"), payloadJSON: `{"version":4}`, method: "POST", path: base + "artifacts/tf/", body: `{"version":4}`, wantOut: "Artifact created successfully."},
		{name: "artifact delete", ctor: artifactscmd.NewArtifactsCmd, args: with("delete", "--workflow-id", "w", "--artifact-id", "tf"), method: "DELETE", path: base + "artifacts/tf/", wantOut: "Artifact deleted successfully."},
		{name: "artifact versions", ctor: artifactscmd.NewArtifactsCmd, args: with("versions", "--workflow-id", "w", "--artifact-id", "tf", "--limit", "5"), method: "GET", path: base + "artifacts/tf/versions/", query: url.Values{"limit": {"5"}}},
		{name: "artifact rollback", ctor: artifactscmd.NewArtifactsCmd, args: with("rollback", "--workflow-id", "w", "--artifact-id", "tf", "--version-id", "v"), method: "POST", path: base + "artifacts/tf/rollback/", body: `{"VersionId":"v"}`, wantOut: "Artifact rolled back."},
		{name: "artifact lock", ctor: artifactscmd.NewArtifactsCmd, args: with("lock", "--workflow-id", "w", "--artifact-id", "tf"), method: "POST", path: base + "artifacts/tf/lock/", resp: `{"ID":"x"}`, wantOut: "Artifact locked."},
		{name: "artifact unlock", ctor: artifactscmd.NewArtifactsCmd, args: with("unlock", "--workflow-id", "w", "--artifact-id", "tf"), method: "DELETE", path: base + "artifacts/tf/lock/", resp: `{"ID":"x"}`, wantOut: "Artifact unlocked."},
	})
}

func TestStackCommands(t *testing.T) {
	base := org + "wfgrps/g/stacks/"
	scope := []string{"--org", "demo-org", "--workflow-group", "g"}
	with := func(args ...string) []string { return append(args, scope...) }
	runMatrix(t, []matrixCase{
		{name: "read", ctor: stackcmd.NewStackCmd, args: with("read", "--stack-id", "s"), method: "GET", path: base + "s/"},
		{name: "update", ctor: stackcmd.NewStackCmd, args: with("update", "--stack-id", "s", "@payload"), payloadJSON: `{"Description":"d"}`, method: "PATCH", path: base + "s/", body: `{"Description":"d"}`, wantOut: "Stack updated successfully."},
		{name: "list", ctor: stackcmd.NewStackCmd, args: with("list", "--tags", "t"), method: "GET", path: base + "listall/", query: url.Values{"Tags": {"t"}}},
		{name: "compare", ctor: stackcmd.NewStackCmd, args: with("compare", "--stack-id", "s", "--target-template-group-id", "/o/t", "--upgrade-mode", "MANUAL"), method: "POST", path: base + "s/compare/", body: `{"targetTemplateGroupId":"/o/t","upgradeMode":"MANUAL"}`, wantOut: `"msg"`},
		{name: "run create", ctor: stackrun.NewStackRunCmd, args: with("create", "--stack-id", "s", "--action-type", "APPLY"), method: "POST", path: base + "s/stackruns/", body: `{"ActionType":"APPLY"}`, wantOut: "Stack run created successfully."},
		{name: "run read", ctor: stackrun.NewStackRunCmd, args: with("read", "--stack-id", "s", "--run-id", "r1"), method: "GET", path: base + "s/stackruns/r1/"},
		{name: "run list", ctor: stackrun.NewStackRunCmd, args: with("list", "--stack-id", "s", "--limit", "5"), method: "GET", path: base + "s/stackruns/listall/", query: url.Values{"limit": {"5"}}},
	})
}

func TestStackWorkflowCommands(t *testing.T) {
	base := org + "wfgrps/g/stacks/s/wfs/"
	scope := []string{"--org", "demo-org", "--workflow-group", "g", "--stack-id", "s"}
	with := func(args ...string) []string { return append(args, scope...) }
	runMatrix(t, []matrixCase{
		{name: "create", ctor: stackworkflow.NewStackWorkflowCmd, args: with("create", "@payload"), payloadJSON: `{"ResourceName":"w"}`, method: "POST", path: base, body: `{"ResourceName":"w"}`, wantOut: "Stack workflow created successfully."},
		{name: "read", ctor: stackworkflow.NewStackWorkflowCmd, args: with("read", "--workflow-id", "w"), method: "GET", path: base + "w/"},
		{name: "update", ctor: stackworkflow.NewStackWorkflowCmd, args: with("update", "--workflow-id", "w", "@payload"), payloadJSON: `{"Description":"d"}`, method: "PATCH", path: base + "w/", body: `{"Description":"d"}`, wantOut: "Stack workflow updated successfully."},
		{name: "delete", ctor: stackworkflow.NewStackWorkflowCmd, args: with("delete", "--workflow-id", "w"), method: "DELETE", path: base + "w/", wantOut: "Stack workflow deleted successfully."},
		{name: "list", ctor: stackworkflow.NewStackWorkflowCmd, args: with("list", "--latest-wf-run-statuses", "3", "--limit", "5"), method: "GET", path: base + "listall/", query: url.Values{"LatestWfRunStatuses": {"3"}, "limit": {"5"}}},
		{name: "outputs", ctor: stackworkflow.NewStackWorkflowCmd, args: with("outputs", "--workflow-id", "w"), method: "GET", path: base + "w/outputs/"},
		{name: "tfstate url", ctor: stackworkflow.NewStackWorkflowCmd, args: with("tfstate-upload-url", "--workflow-id", "w", "--filename", "s.json"), method: "GET", path: base + "w/tfstate_upload_url/", query: url.Values{"filename": {"s.json"}}},
		{name: "file upload url", ctor: stackworkflow.NewStackWorkflowCmd, args: with("file-upload-url", "--workflow-id", "w", "--filename", "f"), method: "GET", path: base + "w/file_upload_url/", query: url.Values{"filename": {"f"}}},
		{name: "compare", ctor: stackworkflow.NewStackWorkflowCmd, args: with("compare", "--workflow-id", "w", "--target-template-id", "/o/t"), method: "POST", path: base + "w/compare/", body: `{"targetTemplateId":"/o/t"}`, wantOut: `"msg"`},
		{name: "vcs triggers", ctor: stackworkflow.NewStackWorkflowCmd, args: with("vcs-triggers", "--workflow-id", "w", "@payload"), payloadJSON: `{}`, method: "POST", path: base + "w/webhooks/vcs_triggers/", body: `{"VCSConfig":null,"VCSTriggers":null}`, wantOut: "VCS triggers created successfully."},
		{name: "artifacts list", ctor: stackworkflow.NewStackWorkflowCmd, args: with("artifacts", "list", "--workflow-id", "w"), method: "GET", path: base + "w/listall_artifacts/"},
		{name: "artifact url", ctor: stackworkflow.NewStackWorkflowCmd, args: with("artifacts", "get-url", "--workflow-id", "w", "--artifact-path", "p"), method: "GET", path: base + "w/get_artifact/", query: url.Values{"artifactPath": {"p"}}},
		{name: "artifact read", ctor: stackworkflow.NewStackWorkflowCmd, args: with("artifacts", "read", "--workflow-id", "w", "--artifact-id", "tf"), method: "GET", path: base + "w/artifacts/tf/", resp: `{"v":1}`},
		{name: "artifact create", ctor: stackworkflow.NewStackWorkflowCmd, args: with("artifacts", "create", "--workflow-id", "w", "--artifact-id", "tf", "@payload"), payloadJSON: `[1,2]`, method: "POST", path: base + "w/artifacts/tf/", body: `[1,2]`, wantOut: "Artifact created successfully."},
		{name: "artifact delete", ctor: stackworkflow.NewStackWorkflowCmd, args: with("artifacts", "delete", "--workflow-id", "w", "--artifact-id", "tf"), method: "DELETE", path: base + "w/artifacts/tf/", wantOut: "Artifact deleted successfully."},
		{name: "artifact lock", ctor: stackworkflow.NewStackWorkflowCmd, args: with("artifacts", "lock", "--workflow-id", "w", "--artifact-id", "tf"), method: "POST", path: base + "w/artifacts/tf/lock/", resp: `{"ID":"x"}`, wantOut: "Artifact locked."},
		{name: "artifact unlock", ctor: stackworkflow.NewStackWorkflowCmd, args: with("artifacts", "unlock", "--workflow-id", "w", "--artifact-id", "tf"), method: "DELETE", path: base + "w/artifacts/tf/lock/", resp: `{"ID":"x"}`, wantOut: "Artifact unlocked."},
	})
}

func TestRunCommands(t *testing.T) {
	base := org + "wfgrps/g/wfs/w/wfruns/"
	sbase := org + "wfgrps/g/stacks/s/wfs/w/wfruns/"
	scope := []string{"--org", "demo-org", "--workflow-group", "g", "--workflow-id", "w"}
	with := func(args ...string) []string { return append(args, scope...) }
	sscope := []string{"--org", "demo-org", "--workflow-group", "g", "--stack-id", "s", "--workflow-id", "w"}
	swith := func(args ...string) []string { return append(args, sscope...) }
	runMatrix(t, []matrixCase{
		{name: "create empty", ctor: workflowrun.NewWorkflowRunCmd, args: with("create"), method: "POST", path: base, body: `{}`, wantOut: "Workflow run created successfully."},
		{name: "create payload", ctor: workflowrun.NewWorkflowRunCmd, args: with("create", "@payload"), payloadJSON: `{"WfType":"TERRAFORM"}`, method: "POST", path: base, body: `{"WfType":"TERRAFORM"}`, wantOut: "Workflow run created successfully."},
		{name: "read", ctor: workflowrun.NewWorkflowRunCmd, args: with("read", "--run-id", "r1"), method: "GET", path: base + "r1/"},
		{name: "update", ctor: workflowrun.NewWorkflowRunCmd, args: with("update", "--run-id", "r1", "@payload"), payloadJSON: `{"WfType":"CUSTOM"}`, method: "PATCH", path: base + "r1/", body: `{"WfType":"CUSTOM"}`, wantOut: "Workflow run updated successfully."},
		{name: "cancel", ctor: workflowrun.NewWorkflowRunCmd, args: with("cancel", "--run-id", "r1"), method: "PATCH", path: base + "r1/cancel/", wantOut: "Workflow run cancelled."},
		{name: "delete", ctor: workflowrun.NewWorkflowRunCmd, args: with("delete", "--run-id", "r1"), method: "DELETE", path: base + "r1/", wantOut: "Workflow run cancelled."},
		{name: "logs", ctor: workflowrun.NewWorkflowRunCmd, args: with("logs", "--run-id", "r1"), method: "GET", path: base + "r1/logs/"},
		{name: "approve", ctor: workflowrun.NewWorkflowRunCmd, args: with("approve", "--run-id", "r1", "--message", "ok"), method: "POST", path: base + "r1/resume/", body: `{"Approve":true,"Message":"ok"}`, wantOut: "Approval decision submitted."},
		{name: "reject", ctor: workflowrun.NewWorkflowRunCmd, args: with("approve", "--run-id", "r1", "--reject", "--approval-step", "plan"), method: "POST", path: base + "r1/resume/", body: `{"Approve":false,"ApprovalStep":"plan"}`, wantOut: "Approval decision submitted."},
		{name: "list", ctor: workflowrun.NewWorkflowRunCmd, args: with("list", "--context-tags", "c", "--limit", "5"), method: "GET", path: base + "listall/", query: url.Values{"contextTags": {"c"}, "limit": {"5"}}},
		{name: "read by ksuid", ctor: workflowrun.NewWorkflowRunCmd, args: []string{"read-by-ksuid", "--org", "demo-org", "--parent-ksuid", "pk", "--resource-ksuid", "rk"}, method: "GET", path: org + "wfs/pk/wfruns/rk/"},
		{name: "update by ksuid", ctor: workflowrun.NewWorkflowRunCmd, args: []string{"update-by-ksuid", "--org", "demo-org", "--parent-ksuid", "pk", "--resource-ksuid", "rk", "@payload"}, payloadJSON: `{"WfType":"CUSTOM"}`, method: "PATCH", path: org + "wfs/pk/wfruns/rk/", body: `{"WfType":"CUSTOM"}`, wantOut: "Workflow run updated successfully."},
		{name: "facts read", ctor: workflowrun.NewWorkflowRunCmd, args: with("facts", "read", "--run-id", "r1", "--facts-id", "f1"), method: "GET", path: base + "r1/wfrunfacts/f1/", resp: `{"facts":1}`},
		{name: "facts create", ctor: workflowrun.NewWorkflowRunCmd, args: with("facts", "create", "--run-id", "r1", "--facts-id", "f1", "@payload"), payloadJSON: `{"TfDrift":{"a":1}}`, method: "POST", path: base + "r1/wfrunfacts/f1/", body: `{"TfDrift":{"a":1}}`, resp: `{"msg":"ok"}`, wantOut: "Workflow run facts created."},
		{name: "facts update", ctor: workflowrun.NewWorkflowRunCmd, args: with("facts", "update", "--run-id", "r1", "--facts-id", "f1", "@payload"), payloadJSON: `{"Tags":["a"]}`, method: "PATCH", path: base + "r1/wfrunfacts/f1/", body: `{"Tags":["a"]}`, resp: `{"msg":"ok"}`, wantOut: "Workflow run facts updated."},
		{name: "stack create empty", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("create"), method: "POST", path: sbase, body: `{}`, wantOut: "Stack workflow run created successfully."},
		{name: "stack read", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("read", "--run-id", "r1"), method: "GET", path: sbase + "r1/"},
		{name: "stack update", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("update", "--run-id", "r1", "@payload"), payloadJSON: `{"WfType":"CUSTOM"}`, method: "PATCH", path: sbase + "r1/", body: `{"WfType":"CUSTOM"}`, wantOut: "Stack workflow run updated successfully."},
		{name: "stack delete", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("delete", "--run-id", "r1"), method: "DELETE", path: sbase + "r1/", wantOut: "Stack workflow run cancelled."},
		{name: "stack logs", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("logs", "--run-id", "r1"), method: "GET", path: sbase + "r1/logs/"},
		{name: "stack approve", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("approve", "--run-id", "r1", "--reject"), method: "POST", path: sbase + "r1/resume/", body: `{"Approve":false}`, wantOut: "Approval decision submitted."},
		{name: "stack list", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("list", "--latest-status", "RUNNING", "--limit", "5"), method: "GET", path: sbase + "listall/", query: url.Values{"latestStatus": {"RUNNING"}, "limit": {"5"}}, resp: `{"msg":[]}`},
		{name: "stack facts read", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("facts", "read", "--run-id", "r1", "--facts-id", "f1"), method: "GET", path: sbase + "r1/wfrunfacts/f1/", resp: `{"facts":1}`},
		{name: "stack facts create", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("facts", "create", "--run-id", "r1", "--facts-id", "f1", "@payload"), payloadJSON: `{"TfDrift":{}}`, method: "POST", path: sbase + "r1/wfrunfacts/f1/", body: `{"TfDrift":{}}`, resp: `{"msg":"ok"}`, wantOut: "Stack workflow run facts created."},
		{name: "stack facts update", ctor: stackworkflowrun.NewStackWorkflowRunCmd, args: swith("facts", "update", "--run-id", "r1", "--facts-id", "f1", "@payload"), payloadJSON: `{"IsActive":"0"}`, method: "PATCH", path: sbase + "r1/wfrunfacts/f1/", body: `{"IsActive":"0"}`, resp: `{"msg":"ok"}`, wantOut: "Stack workflow run facts updated."},
	})
}

func TestTemplateCommands(t *testing.T) {
	tt := "/api/v1/templatetypes/"
	runMatrix(t, []matrixCase{
		{name: "list", ctor: template.NewTemplateCmd, args: []string{"list", "--org", "demo-org", "--template-type", "IAC", "--owner-orgs", "o1", "--search-query", "q"}, method: "GET", path: tt + "IAC/templates/listall/", query: url.Values{"OwnerOrgs": {"o1"}, "SearchQuery": {"q"}}},
		{name: "list by owner", ctor: template.NewTemplateCmd, args: []string{"list-by-owner", "--org", "demo-org", "--owner-org", "owner-org", "--shared"}, method: "GET", path: "/api/v1/orgs/owner-org/templates/listall/", query: url.Values{"isSharedTemplate": {"true"}}},
		{name: "create", ctor: template.NewTemplateCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"TemplateType":"IAC","TemplateName":"t1"}`, method: "POST", path: "/api/v1/templates/", body: `{"TemplateType":"IAC","TemplateName":"t1","OwnerOrg":"","SourceConfigKind":""}`, resp: `{"TemplateType":"IAC","TemplateName":"t1"}`, wantOut: "Template created successfully."},
		{name: "read", ctor: template.NewTemplateCmd, args: []string{"read", "--org", "demo-org", "--template-id", "t1:5", "--template-type", "IAC"}, method: "GET", path: tt + "IAC/demo-org/t1:5/"},
		{name: "read other owner", ctor: template.NewTemplateCmd, args: []string{"read", "--org", "demo-org", "--template-id", "t1", "--template-type", "IAC_GROUP", "--owner-org", "stackguardian"}, method: "GET", path: tt + "IAC_GROUP/stackguardian/t1/"},
		{name: "update", ctor: template.NewTemplateCmd, args: []string{"update", "--org", "demo-org", "--template-id", "t1:2", "--template-type", "WORKFLOW_STEP", "@payload"}, payloadJSON: `{"TemplateName":"t2"}`, method: "PATCH", path: tt + "WORKFLOW_STEP/demo-org/t1:2/", body: `{"TemplateName":"t2"}`, wantOut: "Template updated successfully."},
		{name: "delete", ctor: template.NewTemplateCmd, args: []string{"delete", "--org", "demo-org", "--template-id", "t1", "--template-type", "IAC_POLICY"}, method: "DELETE", path: tt + "IAC_POLICY/demo-org/t1/", wantOut: "Template deleted successfully."},
		{name: "group iac", ctor: template.NewTemplateCmd, args: []string{"group-iac", "--org", "demo-org", "--template-id", "grp", "--iac-template-id", "sub:1", "--owner-org", "stackguardian"}, method: "GET", path: tt + "IAC_GROUP/stackguardian/grp/IAC/sub:1/"},
		{name: "input schema", ctor: template.NewTemplateCmd, args: []string{"input-schema", "--org", "demo-org", "--template-id", "t1:2", "--template-type", "IAC", "--schema-type", "FORM_JSONSCHEMA"}, method: "GET", path: tt + "IAC/demo-org/t1:2/get_input_schema/", query: url.Values{"SchemaType": {"FORM_JSONSCHEMA"}}},
		{name: "vcs triggers", ctor: template.NewTemplateCmd, args: []string{"vcs-triggers", "--org", "demo-org", "--template-id", "t1", "--template-type", "IAC", "@payload"}, payloadJSON: `{}`, method: "POST", path: tt + "IAC/demo-org/t1/webhooks/vcs_triggers/", body: `{"VCSConfig":null,"VCSTriggers":null}`, wantOut: "VCS triggers created successfully."},
		{name: "subscriptions read", ctor: template.NewTemplateCmd, args: []string{"subscriptions", "read", "--org", "demo-org", "--subscription-type", "IACSubscriptions"}, method: "GET", path: org + "subscriptions/default/", query: url.Values{"subscriptionType": {"IACSubscriptions"}}},
		{name: "subscriptions create", ctor: template.NewTemplateCmd, args: []string{"subscriptions", "create", "--org", "demo-org", "@payload"}, payloadJSON: `{"IACSubscriptions":{},"WfStepSubscriptions":{},"PolicySubscriptions":{},"IACGroupSubscriptions":{}}`, method: "POST", path: org + "subscriptions/", body: `{"IACSubscriptions":{},"WfStepSubscriptions":{},"PolicySubscriptions":{},"IACGroupSubscriptions":{}}`, wantOut: "Subscription created successfully."},
		{name: "subscriptions update", ctor: template.NewTemplateCmd, args: []string{"subscriptions", "update", "--org", "demo-org", "--action", "subscribe", "@payload"}, payloadJSON: `{"IACSubscriptions":{"/o/t":{}}}`, method: "PATCH", path: org + "subscriptions/default/subscribe/", body: `{"IACSubscriptions":{"/o/t":{}}}`, wantOut: "Subscription updated successfully."},
		{name: "artifacts list", ctor: template.NewTemplateCmd, args: []string{"artifacts", "list", "--org", "demo-org", "--template-id", "t:1", "--template-type", "IAC", "--prefix", "p"}, method: "GET", path: tt + "IAC/demo-org/t:1/listall_artifacts/", query: url.Values{"artifactPrefix": {"p"}}},
		{name: "artifact download", ctor: template.NewTemplateCmd, args: []string{"artifacts", "download-url", "--org", "demo-org", "--template-id", "t", "--template-type", "IAC", "--artifact-id", "a", "--owner-org", "stackguardian"}, method: "GET", path: tt + "IAC/stackguardian/t/artifacts/a/download/presigned-url/"},
		{name: "artifact upload", ctor: template.NewTemplateCmd, args: []string{"artifacts", "upload-url", "--org", "demo-org", "--template-id", "t", "--template-type", "IAC_GROUP", "--artifact-id", "a"}, method: "GET", path: tt + "IAC_GROUP/demo-org/t/artifacts/a/upload/presigned-url/"},
		{name: "artifact delete", ctor: template.NewTemplateCmd, args: []string{"artifacts", "delete", "--org", "demo-org", "--template-id", "t", "--template-type", "IAC", "--artifact-id", "a"}, method: "DELETE", path: tt + "IAC/demo-org/t/artifacts/a/", wantOut: "Artifact deleted successfully."},
		{name: "public list", ctor: template.NewTemplateCmd, args: []string{"public", "list", "--org", "demo-org", "--template-type", "IAC", "--owner-orgs", "stackguardian"}, method: "GET", path: "/api/v1/public/templatetypes/IAC/templates/listall/", query: url.Values{"OwnerOrgs": {"stackguardian"}}},
		{name: "public list all", ctor: template.NewTemplateCmd, args: []string{"public", "list-all", "--org", "demo-org", "--template-types", "IAC,IAC_GROUP"}, method: "GET", path: "/api/v1/public/templates/listall/", query: url.Values{"templateTypes": {"IAC,IAC_GROUP"}}},
		{name: "public read", ctor: template.NewTemplateCmd, args: []string{"public", "read", "--org", "demo-org", "--template-id", "t:1", "--template-type", "IAC", "--owner-org", "stackguardian"}, method: "GET", path: "/api/v1/public/templatetypes/IAC/stackguardian/t:1/"},
		// typed template packages
		{name: "wf tpl read", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"read", "--org", "demo-org", "--template-id", "t1"}, method: "GET", path: tt + "IAC/demo-org/t1/"},
		{name: "wf tpl create", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"create", "--org", "demo-org", "--create-first-revision=false", "@payload"}, payloadJSON: `{"TemplateName":"t1"}`, method: "POST", path: "/api/v1/templates/", query: url.Values{"createFirstRevision": {"false"}}, body: `{"TemplateType":"IAC","TemplateName":"t1"}`, wantOut: "Workflow template created successfully."},
		{name: "wf tpl update", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"update", "--org", "demo-org", "--template-id", "t1", "@payload"}, payloadJSON: `{"TemplateName":"t2"}`, method: "PATCH", path: tt + "IAC/demo-org/t1/", body: `{"TemplateName":"t2"}`, wantOut: "Workflow template updated successfully."},
		{name: "wf tpl delete", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"delete", "--org", "demo-org", "--template-id", "t1"}, method: "DELETE", path: tt + "IAC/demo-org/t1/", wantOut: "Workflow template deleted successfully."},
		{name: "wf rev create", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"revisions", "create", "--org", "demo-org", "--template-id", "t1", "@payload"}, payloadJSON: `{"LongDescription":"d"}`, method: "POST", path: tt + "IAC/demo-org/t1/revisions/", body: `{"TemplateType":"IAC","LongDescription":"d"}`, wantOut: "Revision created successfully."},
		{name: "wf rev read", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"revisions", "read", "--org", "demo-org", "--revision-id", "t1:1"}, method: "GET", path: tt + "IAC/demo-org/t1:1/"},
		{name: "wf rev update", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"revisions", "update", "--org", "demo-org", "--revision-id", "t1:1", "@payload"}, payloadJSON: `{"Alias":"a"}`, method: "PATCH", path: tt + "IAC/demo-org/t1:1/", body: `{"Alias":"a"}`, wantOut: "Revision updated successfully."},
		{name: "wf rev delete", ctor: workflowtemplate.NewWorkflowTemplateCmd, args: []string{"revisions", "delete", "--org", "demo-org", "--revision-id", "t1:1", "--keep-parent-template"}, method: "DELETE", path: tt + "IAC/demo-org/t1:1/", query: url.Values{"keepParentTemplate": {"true"}}, wantOut: "Revision deleted successfully."},
		{name: "step tpl read", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"read", "--org", "demo-org", "--template-id", "t1"}, method: "GET", path: tt + "WORKFLOW_STEP/demo-org/t1/"},
		{name: "step tpl create", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"TemplateName":"t1"}`, method: "POST", path: "/api/v1/templates/", body: `{"TemplateType":"WORKFLOW_STEP","TemplateName":"t1"}`, wantOut: "Workflow step template created successfully."},
		{name: "step tpl update", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"update", "--org", "demo-org", "--template-id", "t1", "@payload"}, payloadJSON: `{"ShortDescription":"d"}`, method: "PATCH", path: tt + "WORKFLOW_STEP/demo-org/t1/", body: `{"ShortDescription":"d"}`, wantOut: "Workflow step template updated successfully."},
		{name: "step tpl delete", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"delete", "--org", "demo-org", "--template-id", "t1"}, method: "DELETE", path: tt + "WORKFLOW_STEP/demo-org/t1/", wantOut: "Workflow step template deleted successfully."},
		{name: "step rev create", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"revisions", "create", "--org", "demo-org", "--template-id", "t1", "@payload"}, payloadJSON: `{"LongDescription":"d"}`, method: "POST", path: tt + "WORKFLOW_STEP/demo-org/t1/revisions/", body: `{"TemplateType":"WORKFLOW_STEP","LongDescription":"d"}`, wantOut: "Revision created successfully."},
		{name: "step rev read", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"revisions", "read", "--org", "demo-org", "--revision-id", "t1:1"}, method: "GET", path: tt + "WORKFLOW_STEP/demo-org/t1:1/"},
		{name: "step rev update", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"revisions", "update", "--org", "demo-org", "--revision-id", "t1:1", "@payload"}, payloadJSON: `{"LongDescription":"d"}`, method: "PATCH", path: tt + "WORKFLOW_STEP/demo-org/t1:1/", body: `{"LongDescription":"d"}`, wantOut: "Revision updated successfully."},
		{name: "step rev delete", ctor: workflowsteptemplate.NewWorkflowStepTemplateCmd, args: []string{"revisions", "delete", "--org", "demo-org", "--revision-id", "t1:1"}, method: "DELETE", path: tt + "WORKFLOW_STEP/demo-org/t1:1/", wantOut: "Revision deleted successfully."},
		{name: "stack tpl read", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"read", "--org", "demo-org", "--template-id", "t1"}, method: "GET", path: tt + "IAC_GROUP/demo-org/t1/"},
		{name: "stack tpl create", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"create", "--org", "demo-org", "@payload"}, payloadJSON: `{"OwnerOrg":"demo-org"}`, method: "POST", path: "/api/v1/templates/", body: `{"TemplateType":"IAC_GROUP","OwnerOrg":"demo-org"}`, wantOut: "Stack template created successfully."},
		{name: "stack tpl update", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"update", "--org", "demo-org", "--template-id", "t1", "@payload"}, payloadJSON: `{"Tags":["a"]}`, method: "PATCH", path: tt + "IAC_GROUP/demo-org/t1/", body: `{"Tags":["a"]}`, wantOut: "Stack template updated successfully."},
		{name: "stack tpl delete", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"delete", "--org", "demo-org", "--template-id", "t1"}, method: "DELETE", path: tt + "IAC_GROUP/demo-org/t1/", wantOut: "Stack template deleted successfully."},
		{name: "stack rev create", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"revisions", "create", "--org", "demo-org", "--template-id", "t1", "@payload"}, payloadJSON: `{"Tags":["a"]}`, method: "POST", path: tt + "IAC_GROUP/demo-org/t1/revisions/", body: `{"TemplateType":"IAC_GROUP","Tags":["a"]}`, wantOut: "Revision created successfully."},
		{name: "stack rev read", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"revisions", "read", "--org", "demo-org", "--revision-id", "t1:1"}, method: "GET", path: tt + "IAC_GROUP/demo-org/t1:1/"},
		{name: "stack rev update", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"revisions", "update", "--org", "demo-org", "--revision-id", "t1:1", "@payload"}, payloadJSON: `{"LongDescription":"d"}`, method: "PATCH", path: tt + "IAC_GROUP/demo-org/t1:1/", body: `{"LongDescription":"d"}`, wantOut: "Revision updated successfully."},
		{name: "stack rev delete", ctor: stacktemplate.NewStackTemplateCmd, args: []string{"revisions", "delete", "--org", "demo-org", "--revision-id", "t1:1", "--keep-parent-template"}, method: "DELETE", path: tt + "IAC_GROUP/demo-org/t1:1/", query: url.Values{"keepParentTemplate": {"true"}}, wantOut: "Revision deleted successfully."},
	})
}
