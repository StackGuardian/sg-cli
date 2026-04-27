## StackGuardian CLI (sg-cli)

> **Note:** This repository hosts release binaries only. The source code is maintained in a private repository. Issues, pull requests, and discussions are not monitored here — please reach out via [support@stackguardian.io](mailto:support@stackguardian.io) or your usual support channel.

---

## Contents

-   [Setup](#setup)
-   [Usage](#usage)
    -   [Required arguments](#required-arguments)
    -   [Optional arguments](#optional-arguments)
    -   [Passing a JSON payload](#passing-a-json-payload)
-   [Examples](#examples)
-   [Git VCS scan and bulk import](#git-vcs-scan-and-bulk-import)
-   [Interactive mode](#interactive-mode)

---

## Setup


Set the following environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `SG_API_TOKEN` | Yes | — | Your StackGuardian API token. Find this in your account settings. |
| `SG_BASE_URL` | No | `https://api.app.stackguardian.io` | StackGuardian API base URL. |

The `sg-cli` also requires [jq](https://jqlang.github.io/jq/download/) for JSON processing. Install it before running any commands.

---

## Usage

The script accepts a JSON payload as its final input. The payload holds information about `ResourceName`, `TemplateConfig`, and so on.

Run the following to see the help menu:

```
./sg-cli stack create

```

### Required arguments

```
--org
--workflow-group

```

### Optional arguments

```
--wait
--run
--preview
--dry-run
--stack-name
--patch-payload

```

### Passing a JSON payload

Pass the JSON payload at the end of all arguments, after `--`. Only one argument is accepted after `--` — providing more will result in an error. All other arguments (required and optional) must be passed before `--`, in any order.

---

## Examples

The examples below use the following base payload:

```json
{
  "ResourceName": "test",
  "TemplatesConfig": {
    "templateGroupId": "/demo-org/azure-stack-demo:1",
    "templates": [
      {
        "id": 0,
        "WfType": "TERRAFORM",
        "ResourceName": "azure33f-vnet-3vXY"
      },
      {
        "id": 1,
        "WfType": "TERRAFORM",
        "ResourceName": "azure_aks-Wngq"
      }
    ]
  }
}

```
<details>
<summary>Example 1: Simple run</summary>

### Hidden Details
This content can be revealed or hidden.
- Supports Markdown
- Supports code blocks

</details>

<details>
<summary>Example 2: Override ResourceName</summary>

```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --resourceName custom_name -- payload.json

```

The payload will be updated:

```json
{
  "ResourceName": "custom_name",
  ...
}

```
</details>

<details>
<summary>Example 3: Patch payload fields</summary>

> Make sure to surround the patch JSON in single quotes `''`, and each key and value with `""`.

```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --patch-payload '{"ResourceName": "custom_name", "TemplatesConfig": {"templates": [{"ResourceName": "first_item"}]}}' -- payload.json

```

The payload will look like the following:

```json
{
  "ResourceName": "custom_name",
  "TemplatesConfig": {
    "templateGroupId": "/demo-org/azure-stack-demo:1",
    "templates": [
      {
        "id": 0,
        "WfType": "TERRAFORM",
        "ResourceName": "first_item"
      },
      {
        "id": 1,
        "WfType": "TERRAFORM",
        "ResourceName": "azure_aks-Wngq"
      }
    ]
  }
}

```
</details>

<details>
<summary>Example 4: Unset an array</summary>

```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --patch-payload '{"TemplatesConfig": {"templates": []}}' -- payload.json

```

> When an array is set to `[]`, it will use the default value.

The payload will look like the following:

```json
{
  "ResourceName": "test",
  "TemplatesConfig": {
    "templateGroupId": "/demo-org/azure-stack-demo:1",
    "templates": []
  }
}

```
</details>

<details>
<summary>Example 5: Add a new key</summary>

```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --patch-payload '{"custom_key": "custom_value"}' -- payload.json

```

> The new key/value will be added to the payload.

The payload will look like the following:

```json
{
  "ResourceName": "test",
  ...
  "custom_key": "custom_value"
}

```
</details>

<details>
<summary>Example 6: Bulk onboard cloud accounts</summary>

```
./sg-cli aws integrate --org demo-org -- payload.json

```

> The payload must contain an array of AWS account objects under the key `awsAccounts`.

```json
{
  "awsAccounts": [
    {
      "ResourceName": "Dummy123",
      "Description": "dummy account",
      "Settings": {
        "kind": "AWS_STATIC",
        "config": [
          {
            "awsAccessKeyId": "hi-its-me-a-dummy-account",
            "awsSecretAccessKey": "keep-your-secrets-safe",
            "awsDefaultRegion": "us-east-1"
          }
        ]
      }
    },
    {
      "ResourceName": "Dummy11345",
      "Description": "dummy account",
      "Settings": {
        "kind": "AWS_STATIC",
        "config": [
          {
            "awsAccessKeyId": "hi-its-me-a-dummy-account",
            "awsSecretAccessKey": "keep-your-secrets-safe",
            "awsDefaultRegion": "us-east-1"
          }
        ]
      }
    }
  ]
}

```
</details>

<details>
<summary>Example 7: Bulk create workflows with tfstate files</summary>

```
./sg-cli workflow create --bulk --org demo-org --workflow-group demo-grp -- payload.json

```

> The payload must contain an array of workflow objects.

```json
[
  {
    "Approvers": [],
    "CLIConfiguration": {
      "TfStateFilePath": "/Users/richie/Documents/StackGuardian/stackguardian-migrator/transformer/tfc/../../out/state-files/aws-terraform.tfstate",
      "WorkflowGroup": {"name":"test2"}
    },
    "DeploymentPlatformConfig": [
      {
        "kind": "AWS_RBAC",
        "config": {
          "integrationId": "/integrations/xyz",
          "profileName": "default"
        }
      }
    ],
    "Description": "",
    "EnvironmentVariables": [],
    "MiniSteps": {
      "notifications": {
        "email": {
          "APPROVAL_REQUIRED": [],
          "CANCELLED": [],
          "COMPLETED": [],
          "ERRORED": []
        }
      },
      "wfChaining": { "COMPLETED": [], "ERRORED": [] }
    },
    "ResourceName": "cli-5",
    "RunnerConstraints": { "type": "shared" },
    "Tags": [],
    "TerraformConfig": {
      "approvalPreApply": false,
      "managedTerraformState": true,
      "terraformVersion": "1.5.3"
    },
    "UserSchedules": [],
    "VCSConfig": {
      "iacInputData": { "data": {}, "schemaType": "RAW_JSON" },
      "iacVCSConfig": {
        "customSource": {
          "config": {
            "auth": "/integrations/github_com",
            "includeSubModule": false,
            "isPrivate": true,
            "ref": "",
            "repo": "https://github.com/joscheuerer/terraform-aws-vpc",
            "workingDir": ""
          },
          "sourceConfigDestKind": "GITHUB_COM"
        },
        "useMarketplaceTemplate": false
      }
    },
    "WfType": "TERRAFORM"
  },
  {
    "Approvers": [],
    "CLIConfiguration": {
      "TfStateFilePath": "/Users/richie/Documents/StackGuardian/stackguardian-migrator/transformer/tfc/../../out/state-files/aws-terraform.tfstate",
      "WorkflowGroup": {"name":"test1"}
    },
    "DeploymentPlatformConfig": [
      {
        "kind": "AWS_RBAC",
        "config": {
          "integrationId": "/integrations/xyz",
          "profileName": "default"
        }
      }
    ],
    "Description": "",
    "EnvironmentVariables": [],
    "MiniSteps": {
      "notifications": {
        "email": {
          "APPROVAL_REQUIRED": [],
          "CANCELLED": [],
          "COMPLETED": [],
          "ERRORED": []
        }
      },
      "wfChaining": { "COMPLETED": [], "ERRORED": [] }
    },
    "ResourceName": "cli-5",
    "RunnerConstraints": { "type": "shared" },
    "Tags": [],
    "TerraformConfig": {
      "approvalPreApply": false,
      "managedTerraformState": true,
      "terraformVersion": "1.5.3"
    },
    "UserSchedules": [],
    "VCSConfig": {
      "iacInputData": { "data": {}, "schemaType": "RAW_JSON" },
      "iacVCSConfig": {
        "customSource": {
          "config": {
            "auth": "/integrations/github_com",
            "includeSubModule": false,
            "isPrivate": true,
            "ref": "",
            "repo": "https://github.com/joscheuerer/terraform-aws-vpc",
            "workingDir": ""
          },
          "sourceConfigDestKind": "GITHUB_COM"
        },
        "useMarketplaceTemplate": false
      }
    },
    "WfType": "TERRAFORM"
  }
]

```
</details>

<details>
<summary>Example 8: Run compliance discovery against integrations</summary>

```
./sg-cli compliance aws --org demo-org --region eu-central-1 --integration-name aws-integ -- payload.json
./sg-cli compliance azure --org demo-org --integration-name aws-integ -- payload.json

```

```json
{
    "VCSConfig": {},
    "WfStepsConfig": [
        {
            "wfStepTemplateId": "/stackguardian/steampipe:2",
            "name": "steampipe",
            "approval": false,
            "timeout": 5400,
            "wfStepInputData": {
                "schemaType": "FORM_JSONSCHEMA",
                "data": {
                    "steampipeCheckArgs": "azure_compliance.benchmark.cis_v150",
                    "awsRegion": "all"
                }
            }
        }
    ],
    "WfType": "CUSTOM"
}

```
</details>

---

## Git VCS scan and bulk import

Scan a GitHub or GitLab organization for Terraform repositories and generate a bulk workflow payload ready for import.

**Step 1: Scan your VCS org**

```bash
# GitHub
./sg-cli git-scan scan --provider github --token ghp_xxx --org my-org

# GitLab
./sg-cli git-scan scan --provider gitlab --token glpat-xxx --org my-group

# With options
./sg-cli git-scan scan --provider github --token ghp_xxx --org my-org \
  --max-repos 50 \
  --wfgrp imported-workflows \
  --vcs-auth /integrations/github_com \
  --output sg-payload.json

```

**CLI options:**

| Flag | Description |
|---|---|
| `--provider`, `-p` | VCS provider: `github` or `gitlab` (required) |
| `--token`, `-t` | VCS access token (required) |
| `--org`, `-o` | GitHub organization or GitLab group |
| `--user`, `-u` | Scan repos for a specific user instead of an org |
| `--max-repos`, `-m` | Maximum repositories to scan (0 = no limit) |
| `--include-archived` | Include archived repositories |
| `--include-forks` | Include forked repositories |
| `--wfgrp` | Workflow group name written into payload (default: `imported-workflows`) |
| `--vcs-auth` | SG VCS integration path (e.g., `/integrations/github_com`) |
| `--managed-state` | Enable SG-managed Terraform state |
| `--output`, `-O` | Output file (default: `sg-payload.json`) |
| `--quiet`, `-q` | Minimal output |
| `--verbose`, `-v` | Debug output |

The scanner detects Terraform directories, infers cloud provider from HCL provider blocks, parses Terraform version from `required_version`, and handles monorepos (each subdirectory becomes a separate workflow).

**Step 2: Review and edit sg-payload.json**

Before importing, fill in the fields the scanner cannot infer automatically:

-   `DeploymentPlatformConfig` — Cloud connector integration ID (AWS/Azure/GCP)
-   `VCSConfig.customSource.config.auth` — VCS integration path for private repos
-   `RunnerConstraints` — `shared` or private runner group

**Step 3: Bulk import to StackGuardian**

```bash
export SG_API_TOKEN=<YOUR_SG_API_TOKEN>
./sg-cli workflow create --bulk --org "<ORG NAME>" -- sg-payload.json

```

---

## Interactive mode

`sg-cli` includes a terminal UI for browsing and managing resources without remembering command syntax.

```bash
./sg-cli interactive
# or
./sg-cli i

```

On launch you will be prompted for your **org** and **workflow group**, which are remembered for the session. From the main menu you can:

-   **Workflows** — list, inspect, and create workflows (single or bulk)
-   **Stacks** — list and inspect stacks
-   **Artifacts** — browse workflow and stack artifacts
-   **Git Scan** — run the VCS scanner interactively
-   **Switch Context** — change org / workflow group mid-session

Navigation: arrow keys to move, Enter to select, Ctrl+C or `q` to go back / exit.

<img width="403" height="305" alt="image" src="https://github.com/user-attachments/assets/da7a48ed-f10a-4c46-be4f-748978db814e" />
