package read

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"testing"

// 	api "github.com/StackGuardian/sg-sdk-go"
// 	"github.com/StackGuardian/sg-sdk-go/client"
// 	option "github.com/StackGuardian/sg-sdk-go/option"
// 	"github.com/stretchr/testify/mock"
// )

// type mockSGSdkClient struct {
// 	mock.Mock
// 	response *api.WorkflowGetResponse
// }

// func (m *mockSGSdkClient) RoundTrip(request *http.Request) (*http.Response, error) {

// 	responseData, err := json.Marshal(m.response)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &http.Response{
// 		Body:       io.NopCloser(bytes.NewReader(responseData)),
// 		Status:     http.StatusText(http.StatusOK),
// 		StatusCode: http.StatusOK,
// 	}, nil
// }

// func TestReadWorkflow(t *testing.T) {
// 	var successExpected api.WorkflowGetResponse
// 	err := json.Unmarshal([]byte(`{
// 	    "msg": {
// 	        "UserJobMemory": 1024.0,
// 	        "UserJobCPU": 512.0,
// 	        "NumberOfApprovalsRequired": 0.0,
// 	        "RunnerConstraints": {
// 	            "type": "shared"
// 	        },
// 	        "IsActive": "1",
// 	        "Approvers": [],
// 	        "Tags": [],
// 	        "DeploymentPlatformConfig": [
// 	            {
// 	                "config": {
// 	                    "profileName": "testAWSConnector",
// 	                    "integrationId": "/integrations/testAWSConnector"
// 	                },
// 	                "kind": "AWS_RBAC"
// 	            }
// 	        ],
// 	        "MiniSteps": {
// 	            "webhooks": {
// 	                "COMPLETED": [
// 	                    {
// 	                        "webhookName": "test",
// 	                        "webhookSecret": "test",
// 	                        "webhookUrl": "test"
// 	                    }
// 	                ],
// 	                "DRIFT_DETECTED": [
// 	                    {
// 	                        "webhookName": "test",
// 	                        "webhookSecret": "test",
// 	                        "webhookUrl": "test"
// 	                    }
// 	                ],
// 	                "ERRORED": [
// 	                    {
// 	                        "webhookName": "test",
// 	                        "webhookSecret": "test",
// 	                        "webhookUrl": "test"
// 	                    }
// 	                ]
// 	            },
// 	            "notifications": {
// 	                "email": {
// 	                    "APPROVAL_REQUIRED": [],
// 	                    "CANCELLED": [],
// 	                    "COMPLETED": [],
// 	                    "ERRORED": []
// 	                }
// 	            },
// 	            "wfChaining": {
// 	                "COMPLETED": [],
// 	                "ERRORED": []
// 	            }
// 	        },
// 	        "Authors": [
// 	            "larisoncarvalho@gmail.com"
// 	        ],
// 	        "WfStepsConfig": [],
// 	        "ActivitySubscribers": [
// 	            "larisoncarvalho@gmail.com"
// 	        ],
// 	        "SubResourceId": "/wfgrps/testWFG/wfs/aws-s3-demo-website-vg6P",
// 	        "OrgId": "/orgs/charming-copper",
// 	        "CreatedAt": 1720772420966.0,
// 	        "IsArchive": "0",
// 	        "Description": "test",
// 	        "ResourceId": "/wfs/aws-s3-demo-website-vg6P",
// 	        "WfType": "TERRAFORM",
// 	        "ModifiedAt": 1721228378694.0,
// 	        "ParentId": "/orgs/charming-copper/wfgrps/testWFG",
// 	        "ResourceType": "WORKFLOW",
// 	        "LatestWfrunStatus": "APPROVAL_REQUIRED",
// 	        "DocVersion": "V3.BETA",
// 	        "EnvironmentVariables": [
// 	            {
// 	                "config": {
// 	                    "textValue": "testvalue",
// 	                    "varName": "test"
// 	                },
// 	                "kind": "PLAIN_TEXT"
// 	            }
// 	        ],
// 	        "EnforcedPolicies": [],
// 	        "ResourceName": "aws-s3-demo-website-vg6P",
// 	        "VCSConfig": {
// 	            "iacVCSConfig": {
// 	                "iacTemplateId": "/stackguardian/aws-s3-demo-website:16",
// 	                "useMarketplaceTemplate": true
// 	            },
// 	            "iacInputData": {
// 	                "schemaType": "FORM_JSONSCHEMA",
// 	                "data": {
// 	                    "bucket_region": "eu-central-1"
// 	                }
// 	            }
// 	        },
// 	        "TerraformConfig": {
// 	            "terraformVersion": "1.5.7",
// 	            "approvalPreApply": true,
// 	            "managedTerraformState": true,
// 	            "terraformPlanOptions": "--run ",
// 	            "postApplyWfStepsConfig": [
// 	                {
// 	                    "name": "post-apply-step-1",
// 	                    "mountPoints": [],
// 	                    "wfStepTemplateId": "/stackguardian/terraform:19",
// 	                    "wfStepInputData": {
// 	                        "schemaType": "FORM_JSONSCHEMA",
// 	                        "data": {
// 	                            "terraformVersion": "1.5.3",
// 	                            "managedTerraformState": true,
// 	                            "terraformAction": "plan-destroy"
// 	                        }
// 	                    },
// 	                    "cmdOverride": "test",
// 	                    "approval": true
// 	                }
// 	            ],
// 	            "prePlanWfStepsConfig": [
// 	                {
// 	                    "name": "pre-plan-step-1",
// 	                    "mountPoints": [],
// 	                    "wfStepTemplateId": "/stackguardian/terraform:19",
// 	                    "wfStepInputData": {
// 	                        "schemaType": "FORM_JSONSCHEMA",
// 	                        "data": {
// 	                            "terraformVersion": "1.4.3",
// 	                            "managedTerraformState": true,
// 	                            "terraformAction": "plan"
// 	                        }
// 	                    },
// 	                    "cmdOverride": "test",
// 	                    "approval": true
// 	                }
// 	            ],
// 	            "preApplyWfStepsConfig": [
// 	                {
// 	                    "name": "pre-apply-step-1",
// 	                    "mountPoints": [],
// 	                    "wfStepTemplateId": "/stackguardian/terraform:19",
// 	                    "wfStepInputData": {
// 	                        "schemaType": "FORM_JSONSCHEMA",
// 	                        "data": {
// 	                            "terraformVersion": "1.4.1",
// 	                            "managedTerraformState": true,
// 	                            "terraformAction": "plan"
// 	                        }
// 	                    },
// 	                    "cmdOverride": "test",
// 	                    "approval": true
// 	                }
// 	            ],
// 	            "driftCheck": true
// 	        }
// 	    }
// 	}`), &successExpected)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// successExpected := api.WorkflowGetResponse{
// 	// 	Msg: &api.GeneratedWorkflowGetResponse{
// 	// 		UserJobMemory:             1024.0,
// 	// 		UserJobCPU:                512.0,
// 	// 		NumberOfApprovalsRequired: 0.0,
// 	// 		RunnerConstraints: &api.RunnerConstraints{
// 	// 			Type: "shared",
// 	// 		},
// 	// 		IsActive: "1",
// 	// 		Approvers: []string{
// 	// 			"larisoncarvalho@gmail.com",
// 	// 		},
// 	// 	},
// 	// }

// 	// successExpected := &api.WorkflowGetResponse{
// 	cases := []struct {
// 		name     string
// 		expected api.WorkflowGetResponse
// 	}{
// 		{
// 			name:     "Success",
// 			expected: successExpected,
// 		},
// 	}

// 	for _, tc := range cases {
// 		mockClient := &mockSGSdkClient{response: &tc.expected}
// 		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
// 		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
// 		cmd := NewReadCmd(c)
// 		cmd.SetArgs([]string{
// 			"--org", "demo-org",
// 			"--workflow-group", "sg-sdk-go-test",
// 			"--workflow-id", "3g9uzt3qksh07bqd5gl7r",
// 		})
// 		b := bytes.NewBufferString("")
// 		cmd.SetOut(b)
// 		cmd.Execute()
// 		out, err := io.ReadAll(b)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		var actualResponse api.WorkflowGetResponse
// 		err = json.Unmarshal(out, &actualResponse)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		// if actualResponse.String() != tc.expected.String() {
// 		// 	t.Fatalf("expected \"%s\" got \"%s\"", "hi", string(out))
// 		// }
// 	}

// }
