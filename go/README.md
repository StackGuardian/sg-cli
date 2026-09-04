# sg-cli (Go)

A command-line client for the StackGuardian platform built on
[`sg-sdk-go`](https://github.com/StackGuardian/sg-sdk-go). Every SDK operation
(226 API endpoints) is reachable from one of the 30 command groups below.

## Contents

- [Build](#build)
- [Configuration](#configuration)
- [Conventions](#conventions)
- [Command groups](#command-groups)
- [Examples](#examples)
- [Payload notes](#payload-notes)
- [Tests](#tests)

## Build

Requires Go 1.22.4 or newer.

```bash
cd go
go build -o sg-cli .
./sg-cli --help
```

## Configuration

| Variable | Required | Default | Description |
|---|---|---|---|
| `SG_API_TOKEN` | Yes | — | API token. Both `sgu_` (personal) and `sgo_` (organization) tokens work; the `apikey ` prefix is added automatically. |
| `SG_BASE_URL` | No | `https://api.app.stackguardian.io` | API base URL. |

## Conventions

- **Organization.** Every group takes `--org <name>` as a persistent flag.
  Groups that work inside a workflow group also take `--workflow-group`, and
  workflow-scoped groups take `--workflow-id`. Nested workflow groups are
  addressed as `parent/child`.
- **Identifiers are flags.** Resources are named with flags such as
  `--stack-id`, `--secret-id`, `--connector-id`, `--template-id`, `--run-id`
  or `--chat-id`. Template revisions use `--revision-id NAME:REVISION`.
- **Payloads.** Commands that send a body take a JSON file as the single
  positional argument after `--`. Pass `-` to read the payload from stdin.
  - `--patch-payload '<json>'` merges keys into the file before sending.
  - `--preview` prints the final payload and continues.
  - `--dry-run` prints the final payload and stops without sending.
- **Output.** Read and list commands print the API response as indented JSON
  on stdout, so they can be piped into `jq`. Write commands print a short
  confirmation; add `--output-json` to print the response as well. List
  endpoints that the API answers with `204 No Content` print `null`.
- **Errors.** API errors are printed as `Error: <status>: <body>` and the
  command exits with status 1.

## Command groups

| Group | Commands |
|---|---|
| `organization` | read, update, create, delete, list, workflows, count-workflows, bulk-action |
| `user` | list, read, create, update, delete |
| `role` | list, read, create, update, delete |
| `role-binding` | read |
| `api-access` | list, read, create, update, delete, regenerate-key |
| `api-token` | create, delete |
| `audit-log` | list |
| `benchmark-report` | read |
| `billing` | details, balance, invoices, dashboard-url, change-plan, `profile setup/update`, `payment-methods list/setup-intent/set-default/detach` |
| `secret` | list, read-bulk, create, update, delete |
| `policy` | list, read, create, update, delete |
| `connector` | list, read, create, update, delete, authenticate, repos, github-repos, list-accounts |
| `connector-group` | list, read, create, update, delete, authenticate, discovery-scan, `connectors list/read/update/delete/authenticate` |
| `runner-group` | list, read, create, update, delete, register, deregister-runner, update-runner-state, storage-backend-auth |
| `state-backend` | list, read, create, update, delete, list-statefiles |
| `resource` | search, move, tags |
| `chat` | list, read, create, update, delete, commit-bundle, create-pr, pr-status, pr-sync-status, `artifacts list/get-url/upload-url`, `messages list/read/create/cancel/retry` |
| `workflow-group` | list, read, create, update, delete, create-child, list-children, list-resources |
| `workflow` | list, read, create, update, delete, apply, destroy, outputs, compare, vcs-triggers, file-upload-url, tfstate-upload-url |
| `workflow-run` | list, read, create, update, delete, cancel, approve, logs, read-by-ksuid, update-by-ksuid, `facts read/create/update` |
| `artifacts` | list, read, create, delete, get-url, versions, rollback, lock, unlock |
| `stack` | list, read, create, update, delete, apply, destroy, outputs, compare |
| `stack-run` | list, read, create |
| `stack-workflow` | list, read, create, update, delete, outputs, compare, vcs-triggers, file-upload-url, tfstate-upload-url, `artifacts list/read/create/delete/get-url/lock/unlock` |
| `stack-workflow-run` | list, read, create, update, delete, approve, logs, `facts read/create/update` |
| `template` | list, list-by-owner, read, create, update, delete, group-iac, input-schema, vcs-triggers, `artifacts list/upload-url/download-url/delete`, `subscriptions read/create/update`, `public list/list-all/read` |
| `workflow-template` | read, create, update, delete, `revisions read/create/update/delete` |
| `workflow-step-template` | read, create, update, delete, `revisions read/create/update/delete` |
| `stack-template` | read, create, update, delete, `revisions read/create/update/delete` |

`sg-cli <group> --help` and `sg-cli <group> <command> --help` list the flags of
each command. Commands that need a Cognito user token rather than an API key
(organization create/list, public templates, template input schemas) say so in
their help text.

## Examples

```bash
export SG_API_TOKEN=sgu_xxx

# Read commands print JSON
./sg-cli organization read --org demo-org
./sg-cli workflow-group list --org demo-org | jq '.msg[].ResourceName'
./sg-cli workflow read --org demo-org --workflow-group demo-grp --workflow-id my-wf
./sg-cli workflow-run logs --org demo-org --workflow-group demo-grp --workflow-id my-wf --run-id abc123
./sg-cli secret read-bulk --org demo-org --names db_password,api_key
./sg-cli template list --org demo-org --template-type IAC
./sg-cli audit-log list --org demo-org --limit 20

# Write commands take a payload after --
./sg-cli secret create --org demo-org -- secret.json
./sg-cli role update --org demo-org --role-id viewer --patch-payload '{"Description":"read only"}' -- role.json
./sg-cli workflow-group create-child --org demo-org --parent platform -- child.json
./sg-cli workflow-group read --org demo-org --workflow-group platform/child

# Preview the merged payload without sending it
./sg-cli connector create --org demo-org --dry-run --patch-payload '{"ResourceName":"aws-prod"}' -- connector.json

# Print the API response of a write command
./sg-cli workflow-template create --org demo-org --output-json -- template.json

# Templates are addressed by NAME, revisions by NAME:REVISION
./sg-cli workflow-template revisions read --org demo-org --revision-id my-template:2
./sg-cli template artifacts list --org demo-org --template-type IAC --template-id my-template:2

# Trigger and approve runs
./sg-cli workflow-run create --org demo-org --workflow-group demo-grp --workflow-id my-wf
./sg-cli workflow-run approve --org demo-org --workflow-group demo-grp --workflow-id my-wf --run-id abc123
```

The `workflow` and `stack` groups keep the original behaviour: `create`
supports `--run`, `--wait`, `--bulk` and state-file upload, as documented in
the [repository README](../README.md).

## Payload notes

These rules come from running each command against a live organization; the
API rejects the request otherwise.

- `policy update` needs `PolicyType` in the payload (`GENERAL` or
  `FILTER.INSIGHT`) even for partial updates.
- `secret update` needs `ResourceName`; the CLI fills it from `--secret-id`.
- Template creates (`workflow-template`, `stack-template`,
  `workflow-step-template`, `template create`) need `TemplateName`,
  `OwnerOrg` in the form `/orgs/<org>` and `SourceConfigKind` (`TERRAFORM`,
  `OPENTOFU`, ... for IaC templates; `MIXED` for stack templates;
  `DOCKER_IMAGE` for step templates). `revisions create` needs the same
  `OwnerOrg` and `SourceConfigKind`.
- A template with revisions cannot be deleted directly; delete its revisions
  first. Deleting the last revision removes the parent.
- `state-backend create` needs `Statefiles` (an empty list is fine) and
  `StateBackendConfig.type` (`aws_s3` or `azure_blob_storage`) with
  `auth.integrationId` pointing at a connector.
- `connector-group create` needs `Settings.awsDefaultRegion` for AWS kinds.
- `api-access create` needs `Roles`.
- `resource search` needs `ResourceTypes` (`--resource-types WORKFLOW,STACK`).
- Runner groups must be enabled for the organization; otherwise the API
  answers `409`.

## Tests

Offline unit tests mock the HTTP transport and check the request each command
sends (method, path, query and body) and what it prints:

```bash
cd go
go test ./tests/ -run 'TestCommandTree|TestMatrix|TestCreateWorkflow|TestReadWorkflow|TestListWorkflow|TestDeleteWorkflow|TestApplyWorkflow|TestDestroyWorkflow|TestUpdateWorkflow|TestCreateStack|TestDeleteStack|TestApplyStack|TestDestroyStack|TestStackOutputs|TestListArtifacts'
```

`tests/tree_test.go` pins the full command tree, so adding or renaming a
command requires updating `expectedLeaves`. The remaining tests in `tests/`
are end-to-end tests that run a built `../sg-cli` binary against a live
organization and need `SG_API_TOKEN`.
