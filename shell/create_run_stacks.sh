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
    --preview           preview payload before applying
    --dry-run           preview payload before applying (but do not create)
    --resource-name     patch payload ResourceName (workflow-stack name) with custom name
    --patch             patch original json payload

  NOTE: --resource-name and --patch can not work together, --patch is higher that --resource-name
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
        --patch)
            readonly json_patch="$2"
            shift 2
            ;;
        --preview)
            readonly preview_patch=true
            shift
            ;;
        --dry-run)
            readonly dry_run=true
            shift
            ;;
        --)
            shift
            if [ $# -gt 1 ]; then
              echo
              echo "ERROR: only file name should be provided after --"
              exit 1
            fi
            payload="$(cat "$1")"
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
    response=$(curl -s --http1.1 -X POST \
      -H 'PrincipalId: ""' \
      -H "Authorization: apikey $api_token" \
      -H "Content-Type: application/json" \
      --data-raw "${payload}" "$url")
    if [ $? -eq 0 ] && echo "$response" | grep -q "\"data\""; then
      echo "$response" | jq
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

    if [ -n "$resource_name" ] && [ -z "$json_patch" ]; then
      payload=$(echo "${payload}" | jq ".ResourceName = \"$resource_name\"")
    elif [ -n "$json_patch" ]; then
      patch_payload "$(get_root_patch_keys)"
    fi
    if [ "${dry_run}" = "true" ]; then
      echo "${payload}" | jq
      exit 0
    elif [ "${preview_patch}" = "true" ]; then
      echo "${payload}" | jq
    fi

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

check_patch_subkey() {
  if echo "${json_patch}" | jq -r "$1 | keys_unsorted[]"; then
    return 0
  fi
  return 1
}

get_root_patch_keys()  {
  echo "${json_patch}" | jq -r "keys_unsorted[]"
}

get_sub_patch_keys() {
  echo "${json_patch}" | jq -r "$1 | keys_unsorted[]"
}

is_sub_patch_key_array() {
  echo "${json_patch}" | jq -r "$1 | if type==\"array\" then \"yes\" else \"no\" end"
}

fetch_patch_array_length() {
  echo "${json_patch}" | jq "$1 | length"
}

patch_payload_array() {
  array_length=$(($(fetch_patch_array_length "$1")-1))
  if [ "${array_length}" -lt 0 ]; then
    payload="$(echo "${payload}" | jq "${1} = []")"
    return 0
  fi
  until [ $array_length -lt 0 ]; do
    patch_payload "$(get_sub_patch_keys "${1}[${array_length}]")" "${1}[${array_length}]"
    array_length=$((array_length-1))
  done
}

patch_payload() {
  root_keys="$1"
  for key in ${root_keys}; do
    if check_patch_subkey "${2}.${key}" >/dev/null 2>&1; then
      if [ "$(is_sub_patch_key_array "${2}.${key}")" = "yes" ]; then
        patch_payload_array "${2}.${key}"
      else
        patch_payload "$(get_sub_patch_keys "${2}.${key}")" "${2}.${key}"
      fi
    else
      payload="$(echo "${payload}" | jq "${2}.${key} = \"$(echo "${json_patch}" | jq -r "${2}.${key}")\"")"
    fi
  done
}

# run main function
main "$@"
