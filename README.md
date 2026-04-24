## StackGuardian CLI (sg-cli)

> **Note:** This repository hosts release binaries only. The source code is
> maintained in a private repository. Issues, pull requests, and discussions
> are not monitored here — please reach out via [support@stackguardian.io](mailto:support@stackguardian.io)
> or your usual support channel.

### 1: Setup

Required environment variables:
```
SG_BASE_URL (default: https://api.app.stackguardian.io)
SG_API_TOKEN
SG_DASHBOARD_URL (default: https://app.stackguardian.io/orchestrator)
```
Install jq in your environment: https://jqlang.github.io/jq/download/

### 2: Required input

Script accepts JSON payload for the final input.
Payload holds information about `ResourceName`, `TemplateConfig` and so on.

### 3: Running script

When running just
```
./sg-cli stack create
```
help menu will be shown with more details.

There are required arguments that need to be passed when running script:
```
--org
--workflow-group
```
and optional like:
```
--wait
--run
--preview
--dry-run
--stack-name
--patch-payload
```
JSON payload is passed at the end of all arguments after `--`.
Only one arguments is accepted after `--`, providing more will result in error.
Any argument (optional, required) needs to be passed before `--`, in any order.

If we have payload like following
```
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

Example 1: (simple run with prefilled payload.json)
```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp -- payload.json
```

Example 2: (override ResourceName (workflow-stack name))
```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --resourceName custom_name -- payload.json

```
Payload from before will have updated:
```
{
  "ResourceName": "custom_name",
  ...
}
```

Example 3: (patch anything inside payload.json)
> make sure to surround patch json in single quotes `''`, and each key and value with `""`
```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --patch-payload '{"ResourceName": "custom_name", "TemplatesConfig": {"templates": [{"ResourceName": "first_item"}]}}' -- payload.json
```
Paylod will look like the following:
```
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

Example 4: (unset array)
```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --patch-payload '{"TemplatesConfig": {"templates": []}}' -- payload.json
```
Payload will look like the follwing:
> when array is set to `[]`, it will use default value
```
{
  "ResourceName": "test",
  "TemplatesConfig": {
    "templateGroupId": "/demo-org/azure-stack-demo:1",
    "templates": []
  }
}
```

Example 5: (add new key)
```
./sg-cli stack create --org demo-org --workflow-group integration-wfgrp --patch-payload '{"custom_key": "custom_value"}' -- payload.json
```
Payload will look like the follwing:
> new key/value will be added to payload
```
{
  "ResourceName": "test",
  ...
  "custom_key": "custom_value"
}
```

Example 6: Bulk onboard cloud accounts
```
./sg-cli aws integrate --org demo-org  -- payload.json
```

Payload will look like the follwing:
> It should contain an array of AWS account objects under the key `awsAccounts`
```
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

Example 7: Bulk create workflows with tfstate files
```
./sg-cli workflow create --bulk --org demo-org --workflow-group demo-grp  -- payload.json
```

payload.json will look like the following:
>  payload.json should contain an array of workflow objects
```
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


Example 8: Run Compliance discovery against integrations
```
./sg-cli compliance aws --org demo-org --region eu-central-1 --integration-name aws-integ -- payload.json
./sg-cli compliance azure --org demo-org --integration-name aws-integ -- payload.json
```

payload.json will look like the following:
>  payload.json example
```
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
    "WfType": "CUSTOM",
}
```

---

### Git VCS Scan + Bulk Import

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

- `DeploymentPlatformConfig` — Cloud connector integration ID (AWS/Azure/GCP)
- `VCSConfig.customSource.config.auth` — VCS integration path for private repos
- `RunnerConstraints` — `shared` or private runner group

**Step 3: Bulk import to StackGuardian**

```bash
export SG_API_TOKEN=<YOUR_SG_API_TOKEN>
./sg-cli workflow create --bulk --org "<ORG NAME>" -- sg-payload.json
```

---

### Interactive Mode

sg-cli includes a terminal UI for browsing and managing resources without remembering command syntax.

```bash
./sg-cli interactive
# or
./sg-cli i
```

On launch you will be prompted for your **org** and **workflow group**, which are remembered for the session. From the main menu you can:

- **Workflows** — list, inspect, and create workflows (single or bulk)
- **Stacks** — list and inspect stacks
- **Artifacts** — browse workflow and stack artifacts
- **Git Scan** — run the VCS scanner interactively
- **Switch Context** — change org / workflow group mid-session

Navigation: arrow keys to move, Enter to select, Ctrl+C or `q` to go back / exit.

<img width="403" height="305" alt="image" src="https://github.com/user-attachments/assets/da7a48ed-f10a-4c46-be4f-748978db814e" />
