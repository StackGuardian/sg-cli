import json
import urllib.request
import urllib.parse
from time import sleep
import os


base_url = os.getenv("SG_BASE_URL", "https://api.app.stackguardian.io/")
api_url = f"{base_url}api/v1"
dashboard_url = os.getenv(
    "SG_DASHBOARD_URL", "https://app.stackguardian.io/orchestrator/"
)

api_token = os.environ.get("SG_API_TOKEN")

org = os.environ.get("SG_ORG")

if not api_token or "sgu_" not in api_token:
    print(
        f'Invalid or no API Token provided. Expecting it in "SG_API_TOKEN" environment variable. Navigate to StackGuardian platform to get your api token: {dashboard_url}orgs/{org}/settings?tab=api_key'
    )
    exit(1)

if not org:
    print('No Org provided. Expecting it in "SG_ORG" environment variable')
    exit(1)


def create_stack(org_id, wfgrp_id, payload, runOnCreate=True):
    url = (
        api_url
        + "/orgs/"
        + org_id
        + "/wfgrps/"
        + wfgrp_id
        + "/stacks/?runOnCreate="
        + str(runOnCreate)  # parametrize runOnCreate
    )
    headers = {
        "PrincipalId": "",
        "Authorization": "apikey " + api_token,
        "Content-Type": "application/json",
    }
    req = urllib.request.Request(
        url=url,
        data=json.dumps(payload).encode("utf-8"),
        headers=headers,
        method="POST",
    )
    try:
        with urllib.request.urlopen(req) as response:
            return response.read().decode()
    except urllib.error.HTTPError as e:
        print("==Stack creation failed==")
        print("url: ", e.url)
        print("status: ", e.status)
        print("message: ", e.read().decode())
        return False


def get_wfruns_in_stackrun(org_id, wfgrp_id, stack_id, stackrun_id):
    url = (
        api_url
        + "/orgs/"
        + org_id
        + "/wfgrps/"
        + wfgrp_id
        + "/stacks/"
        + stack_id
        + stackrun_id
    )
    headers = {
        "PrincipalId": "",
        "Authorization": "apikey " + api_token,
        "Content-Type": "application/json",
    }
    req = urllib.request.Request(
        url=url,
        # data=json.dumps(payload).encode("utf-8"),
        headers=headers,
        method="GET",
    )
    try:
        with urllib.request.urlopen(req) as response:
            return response.read().decode()
    except urllib.error.HTTPError as e:
        print("==Stack creation failed==")
        print("url: ", e.url)
        print("status: ", e.status)
        print("message: ", e.read().decode())
        return False


def get_stack(org_id, wfgrp_id, stack_id):
    url = api_url + "/orgs/" + org_id + "/wfgrps/" + wfgrp_id + "/stacks/" + stack_id
    headers = {
        "PrincipalId": "",
        "Authorization": "apikey " + api_token,
        "Content-Type": "application/json",
    }
    # response = requests.get(url, headers=headers)
    req = urllib.request.Request(url=url, headers=headers, method="GET")
    try:
        with urllib.request.urlopen(req) as response:
            return response.read()
    except urllib.error.HTTPError as e:
        print("url: ", e.url)
        print("status: ", e.status)
        print("message: ", e.read().decode())
        return False


def get_stack_status(org_id, wfgrp_id, stack_id):
    # make an api call to base url using GET method
    response = get_stack(org_id, wfgrp_id, stack_id)
    response = json.loads(response)
    if "msg" in response and "LatestWfStatus" in response["msg"]:
        return response["msg"]["LatestWfStatus"]
    else:
        return False


def get_stackrun_status(org_id, wfgrp_id, stack_id, stackrun_id):
    # make an api call to base url using GET method
    response = get_wfruns_in_stackrun(org_id, wfgrp_id, stack_id, stackrun_id)
    response = json.loads(response)
    if "LatestStatus" in response["msg"]:
        print(response["msg"]["LatestStatus"])
        return response["msg"]["LatestStatus"]
    else:
        return False


def main():
    # create a stack
    org_id = org
    wfgrp_id = "azure-stae3ck-wfs-demo" # parametrize
    payload = """{
    "ResourceName": "test-dedffeeed-j3e33456", # parametrize
    "TemplatesConfig": {
        "templateGroupId": "/demo-org/azure-stack-demo:1",
        "templates": [
            {
                "id": 0,
                "WfType": "TERRAFORM",
                "ResourceName": "azure33f-vnet-3vXY",
                "Description": "",
                "EnvironmentVariables": [],
                "DeploymentPlatformConfig": [
                    {
                        "config": {
                        "integrationId": "/integrations/test-azure"
                        },
                        "kind": "AZURE_STATIC"
                    }
                ],
                "TerraformConfig": {
                    "terraformVersion": "1.3.6",
                    "managedTerraformState": true,
                    "approvalPreApply": false,
                    "driftCheck": false
                },
                "VCSConfig": {
                    "iacVCSConfig": {
                        "useMarketplaceTemplate": true,
                        "iacTemplate": "/demo-org/azure-vnet",
                        "iacTemplateId": "/demo-org/azure-vnet:2"
                    },
                    "iacInputData": {
                        "schemaType": "RAW_JSON",
                        "data": {
                            "vnet_name": "wfdemo", # parametrize but how?
                            "subnet_enforce_private_link_service_network_policies": {},
                            "subnet_prefixes": [
                                "192.16.0.0/24",
                                "192.16.1.0/24",
                                "192.16.2.0/24"
                            ],
                            "resource_group_name": "test",
                            "subnet_delegation": {},
                            "dns_servers": [],
                            "tags": {
                                "ENV": "test"
                            },
                            "subnet_enforce_private_link_endpoint_network_policies": {},
                            "address_space": [
                                "192.16.0.0/16"
                            ],
                            "nsg_ids": {},
                            "route_tables_ids": {},
                            "subnet_service_endpoints": {},
                            "vnet_location": "Germany West Central",
                            "subnet_names": [
                                "subnet1",
                                "subnet2",
                                "subnet3"
                            ]
                        }
                    }
                },
                "MiniSteps": {
                    "wfChaining": {
                        "ERRORED": [],
                        "COMPLETED": []
                    },
                    "notifications": {
                        "email": {
                            "ERRORED": [],
                            "COMPLETED": [],
                            "APPROVAL_REQUIRED": [],
                            "CANCELLED": []
                        }
                    }
                },
                "Approvers": [],
                "GitHubComSync": {
                    "pull_request_opened": {
                        "createWfRun": {
                            "enabled": false
                        }
                    }
                }
            },
            {
                "id": 1,
                "WfType": "TERRAFORM",
                "ResourceName": "azure_aks-Wngq",
                "Description": "",
                "EnvironmentVariables": [],
                "DeploymentPlatformConfig": [
                    {
                        "config": {
                        "integrationId": "/integrations/test-azure"
                        },
                        "kind": "AZURE_STATIC"
                    }
                ],
                "TerraformConfig": {
                    "terraformVersion": "1.3.6",
                    "managedTerraformState": true,
                    "approvalPreApply": false,
                    "driftCheck": false
                },
                "VCSConfig": {
                    "iacVCSConfig": {
                        "useMarketplaceTemplate": true,
                        "iacTemplate": "/demo-org/azure_aks",
                        "iacTemplateId": "/demo-org/azure_aks:1"
                    },
                    "iacInputData": {
                        "schemaType": "FORM_JSONSCHEMA",
                        "data": {
                            "agents_labels": {},
                            "admin_username": "azureuser",
                            "prefix": "seelg",
                            "enable_log_analytics_workspace": false,
                            "resource_group_name": "test",
                            "enable_azure_policy": false,
                            "os_disk_size_gb": 50,
                            "net_profile_outbound_type": "loadBalancer",
                            "rbac_aad_managed": false,
                            "network_plugin": "kubenet",
                            "log_analytics_workspace_sku": "PerGB2018",
                            "enable_role_based_access_control": false,
                            "enable_node_public_ip": false,
                            "agents_tags": {},
                            "identity_type": "SystemAssigned",
                            "cluster_name": "sg-wfs-eedemo-cluster",
                            "vnet_subnet_id": "${workflow::azure-stacfk-wfs-demo.test-ddd-2.azure-vnet-3vXY.outputs.vnet_subnets.value.1}",
                            "log_retention_in_days": 30,
                            "agents_size": "Standard_D2s_v3",
                            "enable_host_encryption": false,
                            "enable_http_application_routing": false,
                            "agents_pool_name": "nodepool",
                            "tags": {},
                            "private_cluster_enabled": false,
                            "agents_count": 2,
                            "enable_kube_dashboard": false,
                            "enable_auto_scaling": false,
                            "sku_tier": "Free",
                            "agents_type": "VirtualMachineScaleSets"
                        }
                    }
                },
                "MiniSteps": {
                    "wfChaining": {
                        "ERRORED": [],
                        "COMPLETED": []
                    },
                    "notifications": {
                        "email": {
                            "ERRORED": [],
                            "COMPLETED": [],
                            "APPROVAL_REQUIRED": [],
                            "CANCELLED": []
                        }
                    }
                },
                "Approvers": [],
                "GitHubComSync": {
                    "pull_request_opened": {
                        "createWfRun": {
                            "enabled": false
                        }
                    }
                }
            },
            {
                "id": 2,
                "WfType": "CUSTOM",
                "ResourceName": "kubernetes-workflow-step-VXLk",
                "Description": "",
                "EnvironmentVariables": [],
                "DeploymentPlatformConfig": [],
                "TerraformConfig": {
                    "terraformVersion": "1.3.6",
                    "managedTerraformState": true,
                    "approvalPreApply": false,
                    "driftCheck": false
                },
                "VCSConfig": {
                    "iacVCSConfig": {
                        "useMarketplaceTemplate": true,
                        "iacTemplate": "/demo-org/kubernetes-workflow-step",
                        "iacTemplateId": "/demo-org/kubernetes-workflow-step:1"
                    },
                    "iacInputData": {
                        "schemaType": "RAW_JSON",
                        "data": {}
                    }
                },
                "WfStepsConfig": [
                    {
                        "wfStepTemplateId": "/demo-org/kubernetes:1",
                        "wfStepTemplate": "/demo-org/kubernetes",
                        "name": "test",
                        "approval": false,
                        "wfStepInputData": {
                            "schemaType": "FORM_JSONSCHEMA",
                            "data": {
                                "kubectlVersion": "1.26.0",
                                "workingNamespace": "default",
                                "cmdToExec": "kubectl",
                                "additionalParameters": "",
                                "kubectlAction": "apply",
                                "dryRun": true,
                                "dryRunOption": "client"
                            }
                        }
                    }
                ],
                "MiniSteps": {
                    "wfChaining": {
                        "ERRORED": [],
                        "COMPLETED": []
                    },
                    "notifications": {
                        "email": {
                            "ERRORED": [],
                            "COMPLETED": [],
                            "APPROVAL_REQUIRED": [],
                            "CANCELLED": []
                        }
                    }
                },
                "Approvers": [],
                "GitHubComSync": {
                    "pull_request_opened": {
                        "createWfRun": {
                            "enabled": false
                        }
                    }
                }
            }
        ]
    }
}"""
    payload = payload if isinstance(payload, dict) else json.loads(payload)
    response = create_stack(
        org_id,
        wfgrp_id,
        payload,
    )
    if response:
        response = json.loads(response)
        print("Stack created")
        print(
            dashboard_url
            + "orgs/"
            + org_id
            + "/wfgrps/"
            + wfgrp_id
            + "/stacks/"
            + payload["ResourceName"]
        )
    else:
        exit(1)

    # get a stack
    stack_id = response.get("data").get("stack").get("ResourceName")
    stack_run_id = response.get("data").get("stack").get("StackRunId")
    while get_stackrun_status(org_id, wfgrp_id, stack_id, stack_run_id) not in [
        "ERRORED",
        "COMPLETED",
        "APPROVAL_REQUIRED",
    ]:
        print("Stack under deployment")
        sleep(5)
    else:
        print(
            f"Stack finished with {get_stack_status(org_id, wfgrp_id, stack_id)} status"
        )
        exit()


if __name__ == "__main__":
    main()
