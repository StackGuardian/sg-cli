## StackGuardian CLI (sg-cli)

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

Integrating AWS Accounts:

```
./sg-cli aws integrate --org demo-org  -- payload.json
```

Payload will look like the following for integrating AWS accounts:
> It should contain an array of AWS account objects under the key `awsAccounts`
```
{
    "awsAccounts": [
        {
            "ResourceName": "AWS-STATIC-101",
            "Description": "Dummy AWS Account integration using Access Key.",
            "Settings": {
                "kind": "AWS_STATIC",
                "config": [
                    {
                        "awsAccessKeyId": "dummy-accesskey-id",
                        "awsSecretAccessKey":"keep-your-secret-safe",
                        "awsDefaultRegion":"us-east-1"
                    }
                ]
            },
            "Tags": [
                "aws",
                "sg-cli",
                "STATIC"
            ]
        },
        {
            "ResourceName": "AWS-RBAC-101",
            "Description": "Dummy AWS Account integration using RBAC.",
            "Settings": {
                "kind": "AWS_RBAC",
                "config": [
                    {
                        "externalId": "demo-org:1234567890",
                        "durationSeconds":"3600",
                        "roleArn":"arn:aws:iam::account-id:role/role-name"
                    }
                ]
            },
            "Tags": [
                "aws",
                "sg-cli",
                "RBAC"
            ]
        }
    ]
}
```

Integrating Azure Subscriptions:

```
./sg-cli azure integrate --org demo-org  -- payload.json
```

Payload will look like the following for integrating Azure Subscriptions:
> It should contain an array of Azure subscriptions objects under the key `azureSubscription`
```
{
    "azureSubscription": [
        {
            "ResourceName": "AZURE-DUMMY-101",
            "Description": "Dummy Azure Account 101.",
            "Settings": {
                "kind": "AZURE_STATIC",
                "config": [
                    {
                        "armClientSecret": "dummy-client-secret101",
                        "armClientId":"dummy-client-id101",
                        "armSubscriptionId":"dummy-subscription-id101",
                        "armTenantId": "dummy-tenant-id101"
                    }
                ]
            },
            "Tags": [
                "azure",
                "sg-cli",
                "integration"
            ]
        },
        {
            "ResourceName": "AZURE-DUMMY-102",
            "Description": "Dummy Azure Account 102.",
            "Settings": {
                "kind": "AZURE_STATIC",
                "config": [
                    {
                        "armClientSecret": "dummy-client-secret102",
                        "armClientId":"dummy-client-id102",
                        "armSubscriptionId":"dummy-subscription-id102",
                        "armTenantId": "dummy-tenant-id102"
                    }
                ]
            },
            "Tags": [
                "azure",
                "sg-cli",
                "integration"
            ]
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