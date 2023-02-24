#!/bin/sh

# Define variables
readonly base_url="${SG_BASE_URL:-"https://api.app.stackguardian.io"}"
readonly api_url="$base_url/api/v1"
readonly dashboard_url="${SG_DASHBOARD_URL:-"https://app.stackguardian.io/orchestrator"}"
readonly api_token="$SG_API_TOKEN"

if ! type jq >/dev/null 2>&1; then
  echo
  echo "ERROR: jq command is required!"
  exit 2
fi

help() {
  cat <<EOF

  /bin/sh $(basename "$0") OPTIONS --org <ORG_NAME> --workflow-group <WF_GROUP_NAME> -- <JSON_PAYLOAD_PATH>

  OPTIONS:
    --wait              wait for stack creation, applicable only when --run is set
    --run               run stack after creation
    --resource-name     patch payload ResourceName with custom name
EOF
}

if [ $# -lt 6 ]; then
  help
  exit 1
fi

# Parse command-line arguments
while [ $# -gt 0 ]; do
    case "$1" in
        --org)
            readonly org="$2"
            shift 2
            ;;
        --workflow-group)
            readonly workflow_group="$2"
            shift 2
            ;;
        --resource-name)
            readonly resource_name="$2"
            shift 2
            ;;
        --wait)
            readonly wait_execution=true
            shift
            ;;
        --run)
            readonly run_on_create=true
            shift
            ;;
        --)
            shift
            if [ $# -gt 1 ]; then
              echo
              echo "ERROR: only file name should be provided after --"
              exit 1
            fi
            readonly payload="$1"
            break
            ;;
        *)
            echo "ERROR: unknown option $1" >&2
            exit 1
            ;;
    esac
done

if [ -z "$api_token" ] || [ "sgu_" != "${api_token:0:4}" ]; then
    echo "Invalid or no API Token provided. Expecting it in \"SG_API_TOKEN\" environment variable. Navigate to StackGuardian platform to get your api token: $dashboard_url/orgs/$org/settings?tab=api_key"
    exit 1
fi

create_stack() {
    org_id=$1
    wfgrp_id=$2
    runOnCreate=${run_on_create:-false}
    url="$api_url/orgs/$org_id/wfgrps/$wfgrp_id/stacks/?runOnCreate=$runOnCreate"
    if [ -n "$resource_name" ]; then
      jq ".ResourceName = \"$resource_name\"" "$payload" > "$payload".new
      mv "$payload".new "$payload" && rm -f "$payload".new
    fi
    response=$(curl -s --http1.1 -X POST \
      -H 'PrincipalId: ""' \
      -H "Authorization: apikey $api_token" \
      -H "Content-Type: application/json" \
      -d @"$payload" "$url")
    if [ $? -eq 0 ] && echo "$response" | grep -q "\"data\""; then
      echo "$response"
    else
      echo "== Stack creation failed =="
      echo "url: $url"
      echo "response: $response"
      exit 1
    fi
}

get_wfruns_in_stackrun() {
    org_id=$1
    wfgrp_id=$2
    stack_id=$3
    stackrun_id=$4
    url="$api_url/orgs/$org_id/wfgrps/$wfgrp_id/stacks/$stack_id/$stackrun_id"
    response=$(curl -s --http1.1 -X GET \
      -H 'PrincipalId: ""' \
      -H "Authorization: apikey $api_token" \
      -H "Content-Type: application/json" \
      "$url")
    if [ $? -ne 0 ] || echo "$response" | grep -q "\"error\""; then
        echo "== Retrieving Workflow Run from StackRun failed =="
        echo "url: $url"
        echo "response: $response"
        exit 1
    fi
    echo "$response"
}

get_stack() {
    org_id=$1
    wfgrp_id=$2
    stack_id=$3
    url="$api_url/orgs/$org_id/wfgrps/$wfgrp_id/stacks/$stack_id"
    response=$(curl -s --http1.1 -X GET \
      -H 'PrincipalId: ""' \
      -H "Authorization: apikey $api_token" \
      -H "Content-Type: application/json" \
      "$url")
    if [ $? -ne 0 ] || echo "$response" | grep -q "\"error\""; then
        echo "== Retrieving Stack failed =="
        echo "url: $url"
        echo "response: $response"
        exit 1
    fi
    echo "$response"
}

get_stack_status() {
  org_id=$1
  wfgrp_id=$2
  stack_id=$3
  response=$(get_stack "$org_id" "$wfgrp_id" "$stack_id")
  if echo "$response" | grep -q '"msg":' && echo "$response" | grep -q '"LatestWfStatus":'; then
    echo "$response" | jq -r '.msg.LatestWfStatus'
  else
    echo false
  fi
}

get_stackrun_status() {
  org_id=$1
  wfgrp_id=$2
  stack_id=$3
  stackrun_id=$4
  response=$(get_wfruns_in_stackrun "$org_id" "$wfgrp_id" "$stack_id" "$stackrun_id")
  if echo "$response" | grep -q '"LatestStatus":'; then
    echo "$response" | jq -r '.msg.LatestStatus'
  else
    echo false
  fi
}

# main function
main() {
    # create a stack
    org_id=$org
    wfgrp_id="$workflow_group"
    response=$(create_stack "$org_id" "$wfgrp_id")
    if [ $? -ne 0 ]; then
      echo "$response"
      exit 1
    fi
    if [ "$response" != "" ]; then
        stack_id=$(echo "$response" | jq -r '.data.stack.ResourceName')
        stack_run_id=$(echo "$response" | jq -r '.data.stack.StackRunId')
        echo "Stack created"
        echo "$dashboard_url/orgs/$org_id/wfgrps/$wfgrp_id/stacks/$stack_id"
    else
        exit 1
    fi

    # check stackrun status
    if [ "$wait_execution" = "true" ] && [ "$run_on_create" = "true" ]; then
      echo "Stack run executed"
      while [ "$(get_stackrun_status "$org_id" "$wfgrp_id" "$stack_id" "$stack_run_id")" != "ERRORED" ] \
          && [ "$(get_stackrun_status "$org_id" "$wfgrp_id" "$stack_id" "$stack_run_id")" != "COMPLETED" ] \
          && [ "$(get_stackrun_status "$org_id" "$wfgrp_id" "$stack_id" "$stack_run_id")" != "APPROVAL_REQUIRED" ]; do
          echo "Stack under deployment..."
          sleep 5
      done
    fi

    # print final stack status
    if [ "$wait_execution" = "true" ] && [ "$run_on_create" = "true" ]; then
      echo "Stack finished with $(get_stack_status "$org_id" "$wfgrp_id" "$stack_id") status"
      exit 0
    else
      echo "Stack created. To run it go to the Dashboard!"
      exit 0
    fi
}

# run main function
main "$@"

