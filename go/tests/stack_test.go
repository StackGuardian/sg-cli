package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	stackcmd "github.com/StackGuardian/sg-cli/cmd/stack"
	api "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	option "github.com/StackGuardian/sg-sdk-go/option"
	"github.com/stretchr/testify/mock"
)

func TestApplyStack(t *testing.T) {
	var successfulStackApplyExpected api.GeneratedStackRunsResponse
	successExpected := []byte(`{
    "msg": "Stack run scheduled",
    "data": {
        "workflowruns": [
            {
                "OrgId": "/orgs/not-an-actual-org",
                "SubResourceId": "/wfgrps/not-an-actual-workflow-group/stacks/Stack-test/wfs/ansible-Rfde/wfruns/gqfxr0tn0rz9",
                "CreatedAt": 1730113178197,
                "ResourceName": "gqfxr0tn0rz9",
                "Authors": [
                    "dummy@dummy.com"
                ],
                "ResourceType": "WORKFLOW_RUN",
                "ModifiedAt": 1730113178197,
                "LatestStatus": "QUEUED",
                "Comments": {
                    "1730113178197": {
                        "comment": "Workflow Run initiated",
                        "createdBy": "dummy@dummy.com"
                    }
                },
                "Statuses": {
                    "pre_0_step": [
                        {
                            "name": "QUEUED",
                            "createdAt": 1730113178197
                        }
                    ]
                },
                "RuntimeParameters": {
                    "tfDriftWfRun": false,
                    "tfDriftIacInputData": {},
                    "miniSteps": {
                        "notifications": {
                            "email": {
                                "APPROVAL_REQUIRED": [],
                                "CANCELLED": [],
                                "COMPLETED": [],
                                "ERRORED": []
                            }
                        },
                        "wfChaining": {
                            "COMPLETED": [],
                            "ERRORED": []
                        }
                    },
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
        ],
        "StackRunId": "/stackruns/ox2fputswzzl"
    }
}`)
	err := json.Unmarshal(successExpected, &successfulStackApplyExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
		statusCode     int
	}{
		{
			name:           "Success",
			expectedString: "To view the Stack run, please visit the following URL:\nhttps://app.stackguardian.io/orchestrator/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group/stacks/not-an-actual-stack?tab=runs\nStack apply executed.\n",
			expectedByte:   successExpected,
			statusCode:     http.StatusOK,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte, statusCode: tc.statusCode}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := stackcmd.NewStackCmd(c)
		cmd.SetArgs([]string{
			"apply",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--stack-id", "not-an-actual-stack",
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

func TestDestroyStack(t *testing.T) {
	var successfulStackDestroyExpected api.GeneratedStackRunsResponse
	successExpected := []byte(`{
    "msg": "Stack run scheduled",
    "data": {
        "workflowruns": [
            {
                "OrgId": "/orgs/not-an-actual-org",
                "SubResourceId": "/wfgrps/not-an-actual-workflow-group/stacks/Stack-test/wfs/ansible-Rfde/wfruns/2utvf01uqo4q",
                "CreatedAt": 1730113178197,
                "ResourceName": "2utvf01uqo4q",
                "Authors": [
                    "dummy@dummy.com"
                ],
                "ResourceType": "WORKFLOW_RUN",
                "ModifiedAt": 1730113178197,
                "LatestStatus": "QUEUED",
                "Comments": {
                    "1730113178197": {
                        "comment": "Workflow Run initiated",
                        "createdBy": "dummy@dummy.com"
                    }
                },
                "Statuses": {
                    "pre_0_step": [
                        {
                            "name": "QUEUED",
                            "createdAt": 1730113178197
                        }
                    ]
                },
                "RuntimeParameters": {
                    "tfDriftWfRun": false,
                    "tfDriftIacInputData": {},
                    "miniSteps": {
                        "notifications": {
                            "email": {
                                "APPROVAL_REQUIRED": [],
                                "CANCELLED": [],
                                "COMPLETED": [],
                                "ERRORED": []
                            }
                        },
                        "wfChaining": {
                            "COMPLETED": [],
                            "ERRORED": []
                        }
                    },
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
        ],
        "StackRunId": "/stackruns/ox2fputswzzl"
    }
}`)
	err := json.Unmarshal(successExpected, &successfulStackDestroyExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
		statusCode     int
	}{
		{
			name:           "Success",
			expectedString: "To view the Stack run, please visit the following URL:\nhttps://app.stackguardian.io/orchestrator/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group/stacks/not-an-actual-stack?tab=runs\nStack Workflow destroy run successfully.\n",
			expectedByte:   successExpected,
			statusCode:     http.StatusOK,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte, statusCode: tc.statusCode}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := stackcmd.NewStackCmd(c)
		cmd.SetArgs([]string{
			"destroy",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--stack-id", "not-an-actual-stack",
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

func TestStackOutput(t *testing.T) {
	var successfulStackOutputsExpected api.GeneratedStackOutputsResponse
	successExpected := []byte(`{
    "msg": "Stack output fetched successfully",
    "data": {
        "/wfs/ansible-0": "signed_url_1",
        "/wfs/3v15qwtk0qz8t4r2z2354": "signed_url_2"
    }
}`)
	err := json.Unmarshal(successExpected, &successfulStackOutputsExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedStruct *api.GeneratedStackOutputsResponse
		expectedByte   []byte
		statusCode     int
	}{
		{
			name:           "Success",
			expectedStruct: &successfulStackOutputsExpected,
			expectedByte:   successExpected,
			statusCode:     http.StatusOK,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte, statusCode: tc.statusCode}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := stackcmd.NewStackCmd(c)
		cmd.SetArgs([]string{
			"outputs",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--stack-id", "not-an-actual-stack",
		})
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		var actualResponse api.GeneratedStackOutputsResponse
		err = json.Unmarshal(out, &actualResponse)
		if err != nil {
			t.Fatal(err)
		}

		if reflect.DeepEqual(actualResponse.Data["/wfs/ansible-0"],
			tc.expectedStruct.Data["/wfs/ansible-0"]) == false {
			t.Fatalf("expected \"%v\" \ngot \"%v\"",
				tc.expectedStruct.Data["/wfs/ansible-0"],
				actualResponse.Data["/wfs/ansible-0"])
		}
		if actualResponse.Msg != tc.expectedStruct.Msg {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedStruct.Msg,
				actualResponse.Msg)
		}

	}
}

func TestNormalCreateStack(t *testing.T) {
	var successStackCreateExpected api.GeneratedStackCreateResponse
	successExpected := []byte(`{
    "msg": "Stack 40y3rfv9xdxwplehcxifj created",
    "data": {
        "stack": {
            "OrgId": "/orgs/not-an-actual-org",
            "SubResourceId": "/wfgrps/not-an-actual-wfg/stacks/40y3rfv9xdxwplehcxifj",
            "CreatedAt": 1730114654217,
            "ResourceName": "40y3rfv9xdxwplehcxifj",
            "Description": "Dummy Stack",
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
            "LatestWfStatus": "UNTRACKED",
            "ResourceType": "STACK",
            "ModifiedAt": 1730114654217,
            "EnvironmentVariables": [
                {
                    "kind": "PLAIN_TEXT",
                    "config": {
                        "varName": "test",
                        "textValue": "testValue"
                    }
                }
            ],
            "DeploymentPlatformConfig": [
                {
                    "kind": "AWS_RBAC",
                    "config": {
                        "integrationId": "/integrations/DummyConnector",
                        "profileName": "DummyConnector"
                    }
                }
            ],
            "UserSchedules": [],
            "TemplatesConfig": {
                "templates": [
                    {
                        "Description": "Dummy Workflow",
                        "Tags": [],
                        "WfType": "CUSTOM",
                        "NumberOfApprovalsRequired": 0,
                        "id": "cc0061e9-a75c-421b-a75b-ef918e9f4b28",
                        "MiniSteps": {
                            "notifications": {
                                "email": {
                                    "APPROVAL_REQUIRED": [],
                                    "CANCELLED": [],
                                    "COMPLETED": [],
                                    "ERRORED": []
                                }
                            },
                            "wfChaining": {
                                "COMPLETED": [],
                                "ERRORED": []
                            }
                        }
                    }
                ],
                "templateGroupId": "/not-an-actual-org/ansible:4"
            },
            "CreationOrder": [
                "/wfs/ansible-0"
            ],
            "DeletionOrder": [
                "/wfs/ansible-0"
            ],
            "Actions": {
                "apply": {
                    "name": "Create",
                    "description": "use this action to create resources in the stack",
                    "default": true,
                    "order": {
                        "/wfs/ansible-0": {
                            "dependencies": [],
                            "parameters": {
                                "TerraformAction": {
                                    "action": "apply"
                                }
                            }
                        }
                    }
                },
                "destroy": {
                    "name": "Destroy",
                    "description": "use this action to destroy resources in the stack",
                    "default": true,
                    "order": {
                        "/wfs/ansible-0": {
                            "dependencies": [],
                            "parameters": {
                                "TerraformAction": {
                                    "action": "destroy"
                                }
                            }
                        }
                    }
                }
            },
            "WorkflowRelationsMap": {
                "cc0061e9-a75c-421b-a75b-ef918e9f4b28": "/wfs/ansible-0"
            }
        },
        "workflows": [
            {
                "OrgId": "/orgs/not-an-actual-org",
                "SubResourceId": "/wfgrps/not-an-actual-wfg/stacks/40y3rfv9xdxwplehcxifj/wfs/5vee1srm13y0mb2ehz429",
                "CreatedAt": 1730114654312,
                "ResourceName": "5vee1srm13y0mb2ehz429",
                "EnforcedPolicies": "Use GET Workflow API",
                "Description": "Dummy Stack",
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
                "VCSConfig": {},
                "WfStepsConfig": [],
                "ResourceType": "WORKFLOW",
                "ModifiedAt": 1730114654312,
                "EnvironmentVariables": [],
                "DeploymentPlatformConfig": [],
                "RunnerConstraints": {
                    "selectors": [
                        "shared"
                    ],
                    "type": "shared"
                },
                "CacheConfig": {},
                "WfType": "CUSTOM",
                "TerraformConfig": {},
                "UserSchedules": [],
                "NumberOfApprovalsRequired": 0
            }
        ]
    }
}`)
	err := json.Unmarshal(successExpected, &successStackCreateExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
		statusCode     int
	}{
		{
			name:           "Success",
			expectedString: "Stack created successfully.\n",
			expectedByte:   successExpected,
			statusCode:     http.StatusOK,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte, statusCode: tc.statusCode}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := stackcmd.NewStackCmd(c)
		cmd.SetArgs([]string{
			"create",
			"--org", "not-an-actual-org",
			"--workflow-group", "not-an-actual-workflow-group",
			"--", filepath.Join(samplePayloadsDir, createStackFile),
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

func TestDeleteStack(t *testing.T) {
	// Success response for a standard stack deletion
	successExpected := []byte(`{
	    "msg": "Stack deleted successfully",
	    "data": {
	        "ResourceName": "test-stack",
	        "ResourceId": "/stacks/test-stack",
	        "ResourceType": "STACK",
	        "OrgId": "/orgs/not-an-actual-org",
	        "SubResourceId": "/wfgrps/not-an-actual-workflow-group/stacks/test-stack",
	        "ParentId": "/orgs/not-an-actual-org/wfgrps/not-an-actual-workflow-group",
	        "CreatedAt": 1730113178197,
	        "ModifiedAt": 1730113178197,
	        "Authors": [
	            "dummy@dummy.com"
	        ],
	        "Description": "Test stack for deletion"
	    }
	}`)

	// Define common test constants
	const (
		orgName         = "not-an-actual-org"
		workflowGroup   = "not-an-actual-workflow-group"
		stackID         = "test-stack"
		nonEmptyStackID = "non-empty-stack"
	)

	cases := []struct {
		name           string
		expectedString string
		expectedByte   []byte
		statusCode     int
		args           []string
		expectError    bool
	}{
		{
			name:           "Success - Standard Delete",
			expectedString: "Stack deleted successfully.\n",
			expectedByte:   successExpected,
			args: []string{
				"delete",
				"--org", orgName,
				"--workflow-group", workflowGroup,
				"--stack-id", stackID,
			},
			expectError: false,
			statusCode:  http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockSGSdkClient{response: tc.expectedByte, statusCode: tc.statusCode}
			mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
			c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
			cmd := stackcmd.NewStackCmd(c)
			cmd.SetArgs(tc.args)
			b := bytes.NewBufferString("")
			cmd.SetOut(b)

			var err error
			if tc.expectError {
				// For error cases, we need to capture the output before it exits
				cmd.SetErr(b)
				err = cmd.Execute()
				// We expect an error, but in our mocked environment it won't actually exit
			} else {
				err = cmd.Execute()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			out, err := io.ReadAll(b)
			if err != nil {
				t.Fatal(err)
			}

			if string(out) != tc.expectedString {
				t.Fatalf("expected \"%s\" got \"%s\"",
					tc.expectedString,
					string(out))
			}
		})
	}
}
