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
```
JSON payload is passed at the end of all arguments after `--`.
Only one arguments is accepted after `--`, providing more will result in error.
Any argument (optional, required) needs to be passed before `--`, in any order.

Example:
```
./create_run_stacks.sh --org demo-org --workflow-group integration-wfgrp --wait -- payload.json
```
