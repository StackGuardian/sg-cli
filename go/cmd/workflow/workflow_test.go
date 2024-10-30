package workflow

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	api "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	option "github.com/StackGuardian/sg-sdk-go/option"
	"github.com/stretchr/testify/mock"
)

type mockSGSdkClient struct {
	mock.Mock
	response []byte
}

func (m *mockSGSdkClient) RoundTrip(request *http.Request) (*http.Response, error) {

	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(m.response)),
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}, nil
}

func TestReadWorkflow(t *testing.T) {
	var successWorkflowReadExpected api.WorkflowGetResponse
	successExpected := []byte(`{
	    "msg": {
          "UserJobMemory": 1024,
          "UserJobCPU": 512,
          "NumberOfApprovalsRequired": 0,
          "IsActive": "1",
          "Authors": [
            "dummy@dummy.com"
          ],
          "ActivitySubscribers": [
            "dummy@dummy.com"
          ],
          "SubResourceId": "/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow",
          "OrgId": "/orgs/not-an-actual-org",
          "CreatedAt": 1720772420966,
          "IsArchive": "0",
          "Description": "test",
          "ResourceId": "/wfs/not-an-actual-workflow",
          "WfType": "TERRAFORM",
          "ModifiedAt": 1721228378694,
          "ParentId": "/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group",
          "ResourceType": "WORKFLOW",
          "LatestWfrunStatus": "APPROVAL_REQUIRED",
          "DocVersion": "V3.BETA",
          "ResourceName": "not-an-actual-workflow",
          "RunnerConstraints": {
            "type": "shared"
          }
        }
	}`)
	err := json.Unmarshal(successExpected, &successWorkflowReadExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedStruct *api.WorkflowGetResponse
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedStruct: &successWorkflowReadExpected,
			expectedByte:   successExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewWorkflowCmd(c)
		cmd.SetArgs([]string{
			"read",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--workflow-id", "not-an-actual-workflow",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		var actualResponse api.WorkflowGetResponse
		err = json.Unmarshal(out, &actualResponse)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(actualResponse.Msg.ActivitySubscribers, tc.expectedStruct.Msg.ActivitySubscribers) == false {
			t.Fatalf("expected \"%v\" \ngot \"%v\"",
				tc.expectedStruct.Msg.ActivitySubscribers,
				actualResponse.Msg.ActivitySubscribers)
		}
		if reflect.DeepEqual(actualResponse.Msg.Authors, tc.expectedStruct.Msg.Authors) == false {
			t.Fatalf("expected \"%v\" \ngot \"%v\"",
				tc.expectedStruct.Msg.Authors,
				actualResponse.Msg.Authors)
		}
		if actualResponse.Msg.Description != tc.expectedStruct.Msg.Description {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedStruct.Msg.Description,
				actualResponse.Msg.Description)
		}
		if actualResponse.Msg.ResourceId != tc.expectedStruct.Msg.ResourceId {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedStruct.Msg.ResourceId,
				actualResponse.Msg.ResourceId)
		}
	}

}

func TestApplyWorkflow(t *testing.T) {
	var successWorkflowApplyExpected api.WorkflowRunCreatePatchResponse
	successExpected := []byte(`{
    "msg": "Workflow Run dispatched",
    "data": {
        "OrgId": "/orgs/not-an-actual-org",
        "SubResourceId": "/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow/wfruns/not-an-actual-workflow-run",
        "CreatedAt": 1729823822294,
        "ResourceName": "not-an-actual-workflow-run",
        "Authors": [
            "dummy@dummy.com"
        ],
        "ResourceType": "WORKFLOW_RUN",
        "ModifiedAt": 1729823822294,
        "LatestStatus": "QUEUED",
        "Comments": {
            "1729823822294": {
                "comment": "Workflow Run initiated",
                "createdBy": "dummy@dummy.com"
            }
        },
        "Statuses": {
            "pre_0_step": [
                {
                    "name": "QUEUED",
                    "createdAt": 1729823822294
                }
            ]
        },
        "RuntimeParameters": {
            "tfDriftWfRun": false,
            "tfDriftIacInputData": {},
            "miniSteps": {},
            "approvers": [],
            "vcsConfig": {
                "iacVCSConfig": {
                    "iacTemplateId": "/not-an-actual-org/ansible-dummy:3",
                    "useMarketplaceTemplate": true
                },
                "iacInputData": {
                    "schemaType": "FORM_JSONSCHEMA",
                    "data": {
                        "bucket_region": "eu-central-1"
                    }
                }
            },
            "vcsTriggers": {},
            "environmentVariables": [
                {
                    "config": {
                        "textValue": "testValue",
                        "varName": "test"
                    },
                    "kind": "PLAIN_TEXT"
                }
            ],
            "wfStepsConfig": [],
            "cacheConfig": {},
            "runnerConstraints": {
                "type": "shared"
            },
            "deploymentPlatformConfig": [
                {
                    "config": {
                        "profileName": "DummyConnector",
                        "integrationId": "/integrations/DummyConnector"
                    },
                    "kind": "AWS_RBAC"
                }
            ],
            "wfType": "CUSTOM",
            "userJobCpu": 512,
            "userJobMemory": 1024,
            "numberOfApprovalsRequired": 0
        },
        "LatestStatusKey": "pre_0_step"
    }
}`)
	err := json.Unmarshal(successExpected, &successWorkflowApplyExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedString: "Workflow apply run successfully.\nTo view the workflow run, please visit the following URL:\nhttps://app.stackguardian.io/orchestrator/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow?tab=runs\n",
			expectedByte:   successExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewWorkflowCmd(c)
		cmd.SetArgs([]string{
			"apply",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--workflow-id", "not-an-actual-workflow",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tc.expectedString {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedString,
				string(out))
		}
	}

}

func TestDestroyWorkflow(t *testing.T) {
	var successWorkflowApplyExpected api.WorkflowRunCreatePatchResponse
	successExpected := []byte(`{
    "msg": "Workflow Run dispatched",
    "data": {
        "OrgId": "/orgs/not-an-actual-org",
        "SubResourceId": "/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow/wfruns/not-an-actual-workflow-run",
        "CreatedAt": 1729823822294,
        "ResourceName": "not-an-actual-workflow-run",
        "Authors": [
            "dummy@dummy.com"
        ],
        "ResourceType": "WORKFLOW_RUN",
        "ModifiedAt": 1729823822294,
        "LatestStatus": "QUEUED",
        "Comments": {
            "1729823822294": {
                "comment": "Workflow Run initiated",
                "createdBy": "dummy@dummy.com"
            }
        },
        "Statuses": {
            "pre_0_step": [
                {
                    "name": "QUEUED",
                    "createdAt": 1729823822294
                }
            ]
        },
        "RuntimeParameters": {
            "tfDriftWfRun": false,
            "tfDriftIacInputData": {},
            "miniSteps": {},
            "approvers": [],
            "vcsConfig": {
                "iacVCSConfig": {
                    "iacTemplateId": "/not-an-actual-org/ansible-dummy:3",
                    "useMarketplaceTemplate": true
                },
                "iacInputData": {
                    "schemaType": "FORM_JSONSCHEMA",
                    "data": {
                        "bucket_region": "eu-central-1"
                    }
                }
            },
            "vcsTriggers": {},
            "environmentVariables": [
                {
                    "config": {
                        "textValue": "testValue",
                        "varName": "test"
                    },
                    "kind": "PLAIN_TEXT"
                }
            ],
            "wfStepsConfig": [],
            "cacheConfig": {},
            "runnerConstraints": {
                "type": "shared"
            },
            "deploymentPlatformConfig": [
                {
                    "config": {
                        "profileName": "DummyConnector",
                        "integrationId": "/integrations/DummyConnector"
                    },
                    "kind": "AWS_RBAC"
                }
            ],
            "wfType": "CUSTOM",
            "userJobCpu": 512,
            "userJobMemory": 1024,
            "numberOfApprovalsRequired": 0
        },
        "LatestStatusKey": "pre_0_step"
    }
}`)
	err := json.Unmarshal(successExpected, &successWorkflowApplyExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedString: "Workflow destroy run successfully.\nTo view the workflow runs, please visit the following URL:\nhttps://app.stackguardian.io/orchestrator/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow?tab=runs\n",
			expectedByte:   successExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewWorkflowCmd(c)
		cmd.SetArgs([]string{
			"destroy",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--workflow-id", "not-an-actual-workflow",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tc.expectedString {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedString,
				string(out))
		}
	}

}

func TestListWorkflow(t *testing.T) {
	var successWorkflowListExpected api.WorkflowsListAll
	successExpected := []byte(`{
    "msg": [
        {
            "GitHubComRepoID": "StackGuardian/template-tf-aws-s3-demo-website",
            "IsActive": "1",
            "Description": "test desc",
            "ResourceId": "/wfs/not-an-actual-workflow1",
            "WfType": "CUSTOM",
            "ModifiedAt": 1729823854332,
            "ParentId": "/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group",
            "LatestWfrunStatus": "UNTRACKED",
            "Tags": [],
            "Authors": [
                "dummy@dummy.com"
            ],
            "ResourceName": "not-an-actual-workflow",
            "SubResourceId": "/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow-1",
            "CreatedAt": 1729823854332
        },
        {
            "GitHubComRepoID": "StackGuardian/template-tf-aws-s3-demo-website",
            "IsActive": "1",
            "Description": "test desc",
            "ResourceId": "/wfs/not-an-actual-workflow2",
            "WfType": "CUSTOM",
            "ModifiedAt": 1729254579019,
            "ParentId": "/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group",
            "LatestWfrunStatus": "ERRORED",
            "Tags": [],
            "Authors": [
                "dummy@dummy.com"
            ],
            "ResourceName": "not-an-actual-workflow-2",
            "SubResourceId": "/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow-2",
            "CreatedAt": 1729254576623
		}
		]
	}
	`)
	err := json.Unmarshal(successExpected, &successWorkflowListExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedString: "> Workflow Name:  not-an-actual-workflow\nDescription:  test desc\n\n> Workflow Name:  not-an-actual-workflow-2\nDescription:  test desc\n\n",
			expectedByte:   successExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewWorkflowCmd(c)
		cmd.SetArgs([]string{
			"list",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tc.expectedString {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedString,
				string(out))
		}
	}

}

func TestDeleteWorkflow(t *testing.T) {
	var successWorkflowDeleteExpected api.GeneratedWorkflowDeleteResponse
	successExpected := []byte(`{
	"msg": "Workflow 5bfani3ieghuferb1xt5z deleted"
	}
	`)

	err := json.Unmarshal(successExpected, &successWorkflowDeleteExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedString: "Workflow deleted successfully.\n",
			expectedByte:   successExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewWorkflowCmd(c)
		cmd.SetArgs([]string{
			"delete",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--workflow-id", "not-an-actual-workflow",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tc.expectedString {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedString,
				string(out))
		}
	}

}

func TestNormalCreateWorkflow(t *testing.T) {
	var successWorkflowCreateExpected api.GeneratedWorkflowCreateResponse
	successExpected := []byte(`{
    "msg": "Workflow not-an-actual-workflow created",
    "data": {
        "OrgId": "/orgs/not-an-actual-org",
        "SubResourceId": "/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow",
        "CreatedAt": 1729823854332,
        "ResourceName": "not-an-actual-workflow",
        "EnforcedPolicies": "Use GET Workflow API",
        "Description": "test desc",
        "Tags": [],
        "Authors": [
            "dummy@dummy.com"
        ],
        "DocVersion": "V3.BETA",
        "IsActive": "1",
        "IsArchive": "0",
        "ActivitySubscribers": [
            "dummy@dummy.com"
        ],
        "LatestWfrunStatus": "UNTRACKED",
        "VCSConfig": {
            "iacVCSConfig": {
                "useMarketplaceTemplate": true,
                "iacTemplateId": "/not-an-actual-org/ansible-dummy:3"
            }
        },
        "WfStepsConfig": [],
        "ResourceType": "WORKFLOW",
        "ModifiedAt": 1729823854332,
        "EnvironmentVariables": [],
        "DeploymentPlatformConfig": [
            {
                "kind": "AWS_RBAC",
                "config": {
                    "profileName": "DummyConnector",
                    "integrationId": "/integrations/DummyConnector"
                }
            }
        ],
        "RunnerConstraints": {
            "type": "shared",
            "selectors": [
                "shared"
            ]
        },
        "CacheConfig": {},
        "WfType": "CUSTOM",
        "TerraformConfig": {},
        "UserSchedules": [],
        "NumberOfApprovalsRequired": 0
    }
}`)
	err := json.Unmarshal(successExpected, &successWorkflowCreateExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedString: "Workflow created successfully.\n",
			expectedByte:   successExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewWorkflowCmd(c)
		cmd.SetArgs([]string{
			"create",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--", "testSamples/create_workflow_normal.json",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tc.expectedString {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedString,
				string(out))
		}
	}

}

func TestBulkCreateWorkflow(t *testing.T) {
	var successWorkflowBulkExpected api.GeneratedWorkflowCreateResponse
	successBulkExpected := []byte(`{
    "msg": "Workflow not-an-actual-workflow created",
    "data": {
        "OrgId": "/orgs/not-an-actual-org",
        "SubResourceId": "/wfgrps/not-an-actual-workflow-group/wfs/not-an-actual-workflow",
        "CreatedAt": 1729823854332,
        "ResourceName": "not-an-actual-workflow",
        "EnforcedPolicies": "Use GET Workflow API",
        "Description": "test desc",
        "Tags": [],
        "Authors": [
            "dummy@dummy.com"
        ],
        "DocVersion": "V3.BETA",
        "IsActive": "1",
        "IsArchive": "0",
        "ActivitySubscribers": [
            "dummy@dummy.com"
        ],
        "LatestWfrunStatus": "UNTRACKED",
        "VCSConfig": {
            "iacVCSConfig": {
                "useMarketplaceTemplate": true,
                "iacTemplateId": "/not-an-actual-org/ansible-dummy:3"
            }
        },
        "WfStepsConfig": [],
        "ResourceType": "WORKFLOW",
        "ModifiedAt": 1729823854332,
        "EnvironmentVariables": [],
        "DeploymentPlatformConfig": [
            {
                "kind": "AWS_RBAC",
                "config": {
                    "profileName": "DummyConnector",
                    "integrationId": "/integrations/DummyConnector"
                }
            }
        ],
        "RunnerConstraints": {
            "type": "shared",
            "selectors": [
                "shared"
            ]
        },
        "CacheConfig": {},
        "WfType": "CUSTOM",
        "TerraformConfig": {},
        "UserSchedules": [],
        "NumberOfApprovalsRequired": 0
    }
}`)
	err := json.Unmarshal(successBulkExpected, &successWorkflowBulkExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedString: ">> Processing workflow: not-an-actual-workflow\nWorkflow created successfully.\n",
			expectedByte:   successBulkExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewWorkflowCmd(c)
		cmd.SetArgs([]string{
			"create",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--bulk",
			"--", "testSamples/create_workflow_bulk.json",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(out) != tc.expectedString {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedString,
				string(out))
		}
	}

}
