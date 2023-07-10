## How to use .sh script to run stack on StackGuardian

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