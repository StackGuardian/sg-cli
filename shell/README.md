## How to use .sh script to run stack on StackGuardian

### 1: Setup

Required environment variables:
```
SG_BASE_URL (default: https://api.app.stackguardian.io)
SG_API_TOKEN
SG_DASHBOARD_URL (default: https://app.stackguardian.io/orchestrator)
```

### 2: Required input

Script accepts JSON payload for the final input.
Payload holds information about `ResourceName`, `TemplateConfig` and so on.

### 3: Running script

When running just
```
./create_run_stacks.sh
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
--resource-name
--patch
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
./create_run_stacks.sh --org demo-org --workflow-group integration-wfgrp -- payload.json
```

Example 2: (override ResourceName (workflow-stack name))
```
./create_run_stacks.sh --org demo-org --workflow-group integration-wfgrp --resourceName custom_name -- payload.json

```
Payload from before will have updated:
```
{
  "ResourceName": "custom_name",
  ...
}
```

Example 3: (patch anything inside payload.json)
> make sure to surround patch json in single quotes '', and each key and value with ""
```
./create_run_stacks.sh --org demo-org --workflow-group integration-wfgrp --patch '{"ResourceName": "custom_name", "TemplatesConfig": {"templates": [{"ResourceName": "first_item"}]}}' -- payload.json
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
