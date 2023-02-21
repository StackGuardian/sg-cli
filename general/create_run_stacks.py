import json
import urllib.request
import urllib.parse
from time import sleep
import os


base_url = os.getenv("SG_BASE_URL", "https://api.app.stackguardian.io/")
api_url = f"{base_url}api/v1"
dashboard_url = f"https://app.stackguardian.io/orchestrator/"

api_token = os.environ.get("SG_API_TOKEN")

org = os.environ.get("SG_ORG")

if not api_token or "sgu_" not in api_token:
    print("Invalid or no API Token provided. Expecting it in \"SG_API_TOKEN\" environment variable. Navigate to StackGuardian platform to get your api token: {dashboard_url}orgs/{org}/settings?tab=api_key")
    exit(1)

if not org:
    print("No Org provided. Expecting it in \"SG_ORG\" environment variable")
    exit(1)

def create_stack(org_id, wfgrp_id, payload, runOnCreate=True):
    url = (
        api_url
        + "/orgs/"
        + org_id
        + "/wfgrps/"
        + wfgrp_id
        + "/stacks/?runOnCreate="
        + str(runOnCreate)
    )
    headers = {
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


def get_stack(org_id, wfgrp_id, stack_id):
    url = api_url + "/orgs/" + org_id + "/wfgrps/" + wfgrp_id + "/stacks/" + stack_id
    headers = {
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


def main():
    # create a stack
    org_id = org
    wfgrp_id = "azure-dev-environment"
    payload = {
        "ResourceName": "Stack-1",
        "Description": "This is a description",
        "TemplatesConfig": {
            "templateGroupId": "/innocent-orange/azure-aks:1",
            "templates": [
                {
                    "id": 0,
                    "WfType": "TERRAFORM",
                    "ResourceName": "azure_virtual_network-GcTz",
                    "Description": "",
                    "EnvironmentVariables": [],
                    "DeploymentPlatformConfig": [
                        {
                            "config": {
                                "integrationId": "/integrations/azure-demo"
                            },
                            "kind": "AZURE_STATIC"
                        }
                    ],
                    "TerraformConfig": {
                        "terraformVersion": "1.3.6",
                        "managedTerraformState": True,
                        "approvalPreApply": False,
                        "driftCheck": False,
                    },
                    "VCSConfig": {
                        "iacVCSConfig": {
                            "useMarketplaceTemplate": True,
                            "iacTemplate": "/stackguardian/azure_virtual_network",
                            "iacTemplateId": "/stackguardian/azure_virtual_network:6",
                        },
                        "iacInputData": {
                            "schemaType": "FORM_JSONSCHEMA",
                            "data": {
                                "address_space": ["10.0.0.0/16"],
                                "dns_servers": [],
                                "nsg_ids": {},
                                "route_tables_ids": {},
                                "subnet_delegation": {},
                                "subnet_enforce_private_link_endpoint_network_policies": {},
                                "subnet_enforce_private_link_service_network_policies": {},
                                "subnet_names": ["subnet1", "subnet2", "subnet3"],
                                "subnet_prefixes": [
                                    "10.0.1.0/24",
                                    "10.0.2.0/24",
                                    "10.0.3.0/24",
                                ],
                                "subnet_service_endpoints": {},
                                "tags": {"ENV": "test"},
                                "use_for_each": False,
                                "vnet_name": "a123cctvenet-demo",
                                "resource_group_name": "test",
                                "vnet_location": "Germany West Central",
                            },
                        },
                    },
                    "MiniSteps": {
                        "wfChaining": {
                            "ERRORED": [],
                            "COMPLETED": [
                                {
                                    "workflowGroupId": "aws-prod-environment",
                                    "workflowId": "azure_kubernetes_cluster_aks-kaC5",
                                }
                            ],
                        },
                        "notifications": {
                            "email": {
                                "ERRORED": [],
                                "COMPLETED": [],
                                "APPROVAL_REQUIRED": [],
                                "CANCELLED": [],
                            }
                        },
                    },
                    "Approvers": [],
                    "GitHubComSync": {
                        "pull_request_opened": {"createWfRun": {"enabled": False}}
                    },
                },
                {
                    "id": 1,
                    "WfType": "TERRAFORM",
                    "ResourceName": "azure_kubernetes_cluster_aks-kaC5",
                    "Description": "",
                    "EnvironmentVariables": [],
                    "DeploymentPlatformConfig": [
                        {
                            "config": {
                                "integrationId": "/integrations/azure-demo"
                            },
                            "kind": "AZURE_STATIC"
                        }
                    ],
                    "TerraformConfig": {
                        "approvalPreApply": True,
                        "terraformVersion": "1.3.6",
                        "driftCheck": False,
                        "managedTerraformState": True,
                    },
                    "VCSConfig": {
                        "iacVCSConfig": {
                            "useMarketplaceTemplate": True,
                            "iacTemplate": "/stackguardian/azure_kubernetes_cluster_aks",
                            "iacTemplateId": "/stackguardian/azure_kubernetes_cluster_aks:3",
                        },
                        "iacInputData": {
                            "schemaType": "FORM_JSONSCHEMA",
                            "data": {
                                "resource_group_name": "test",
                                "cluster_name": "sg-11demo11-cluster",
                                "prefix": "sguior",
                                "client_id": "",
                                "client_secret": "",
                                "admin_username": "azureuser",
                                "agents_size": "Standard_E4bds_v5",
                                "log_analytics_workspace_sku": "PerGB2018",
                                "log_retention_in_days": 30,
                                "agents_count": 2,
                                "public_ssh_key": "",
                                "tags": {},
                                "enable_log_analytics_workspace": True,
                                "vnet_subnet_id": "${workflow::aws-prod-environment.test.azure_virtual_network-GcTz.outputs.vnet_subnets.value.1}",
                                "os_disk_size_gb": 50,
                                "private_cluster_enabled": False,
                                "enable_kube_dashboard": False,
                                "enable_http_application_routing": False,
                                "enable_azure_policy": False,
                                "sku_tier": "Free",
                                "enable_role_based_access_control": False,
                                "rbac_aad_managed": False,
                                "network_plugin": "kubenet",
                                "net_profile_outbound_type": "loadBalancer",
                                "enable_auto_scaling": False,
                                "agents_pool_name": "nodepool",
                                "enable_node_public_ip": False,
                                "agents_labels": {},
                                "agents_type": "VirtualMachineScaleSets",
                                "agents_tags": {},
                                "identity_type": "SystemAssigned",
                                "enable_host_encryption": False,
                            },
                        },
                    },
                    "MiniSteps": {
                        "wfChaining": {
                            "ERRORED": [],
                            "COMPLETED": [
                                {
                                    "workflowGroupId": "aws-prod-environment",
                                    "workflowId": "kubernetes_ingress-SbyW",
                                }
                            ],
                        },
                        "notifications": {
                            "email": {
                                "ERRORED": [],
                                "COMPLETED": [],
                                "APPROVAL_REQUIRED": [],
                                "CANCELLED": [],
                            }
                        },
                    },
                    "Approvers": [],
                    "GitHubComSync": {
                        "pull_request_opened": {"createWfRun": {"enabled": False}}
                    },
                },
                {
                    "id": 2,
                    "WfType": "TERRAFORM",
                    "ResourceName": "kubernetes_ingress-SbyW",
                    "Description": "",
                    "EnvironmentVariables": [],
                    "DeploymentPlatformConfig": [],
                    "TerraformConfig": {
                        "terraformVersion": "1.3.6",
                        "managedTerraformState": True,
                        "approvalPreApply": False,
                        "driftCheck": False,
                    },
                    "VCSConfig": {
                        "iacVCSConfig": {
                            "useMarketplaceTemplate": True,
                            "iacTemplate": "/innocent-orange/kubernetes_ingress",
                            "iacTemplateId": "/innocent-orange/kubernetes_ingress:6",
                        },
                        "iacInputData": {
                            "schemaType": "FORM_JSONSCHEMA",
                            "data": {
                                "client_key": "${workflow::aws-prod-environment.test.azure_kubernetes_cluster_aks-kaC5.outputs.client_key.value}",
                                "resource_group_name": "test",
                                "client_certificate": "${workflow::aws-prod-environment.test.azure_kubernetes_cluster_aks-kaC5.outputs.client_certificate.value}",
                                "cluster_endpoint": "${workflow::aws-prod-environment.test.azure_kubernetes_cluster_aks-kaC5.outputs.host.value}",
                                "cluster_ca_certificate": "${workflow::aws-prod-environment.test.azure_kubernetes_cluster_aks-kaC5.outputs.cluster_ca_certificate.value}",
                            },
                        },
                    },
                    "MiniSteps": {
                        "wfChaining": {"ERRORED": [], "COMPLETED": []},
                        "notifications": {
                            "email": {
                                "ERRORED": [],
                                "COMPLETED": [],
                                "APPROVAL_REQUIRED": [],
                                "CANCELLED": [],
                            }
                        },
                    },
                    "Approvers": [],
                    "GitHubComSync": {
                        "pull_request_opened": {"createWfRun": {"enabled": False}}
                    },
                },
            ],
        },
    }
    response = create_stack(org_id, wfgrp_id, payload)
    if response:
        response = json.loads(response)
        print("CREATE_STACK: ", response)
        print("Stack created")
    else:
        exit(1)

    # get a stack
    stack_id = response.get("data").get("stack").get("ResourceName")
    while get_stack_status(org_id, wfgrp_id, stack_id) not in [
        "ERRORED",
        "COMPLETED",
        "APPROVAL_REQUIRED",
    ]:
        print("Stack under deployment")
        sleep(5)
    else:
        print(f"Stack finished with {get_stack_status(org_id, wfgrp_id, stack_id)} status. Find more info here:")
        print(dashboard_url
        + "orgs/"
        + org_id
        + "/wfgrps/"
        + wfgrp_id
        + "/stacks/"
        + payload["ResourceName"])
        exit()


if __name__ == "__main__":
    main()
