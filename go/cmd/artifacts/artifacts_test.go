package artifacts

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

func TestArtifactsList(t *testing.T) {
	var successfulArtifactsListExpected api.GeneratedWorkflowListAllArtifactsResponse
	successExpected := []byte(`{
    "msg": "Outputs retrieved",
    "data": {
        "artifacts": {
            "orgs/not-an-actual-org/wfgrps/not-an-actual-wfg/wfs/not-an-actual-workflow/artifacts/tfstate.json": {
                "url": "https://<s3-bucket-path>.com/orgs/not-an-actual-org/wfgrps/not-an-actual-wfg/wfs/not-an-actual-workflow/artifacts/tfstate.json",
                "lastModified": "2024-10-17 15:54:00+00:00",
                "size": 6548
            }
        }
    }
}`)
	err := json.Unmarshal(successExpected, &successfulArtifactsListExpected)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name           string
		expectedStruct *api.GeneratedWorkflowListAllArtifactsResponse
		expectedByte   []byte
	}{
		{
			name:           "Success",
			expectedStruct: &successfulArtifactsListExpected,
			expectedByte:   successExpected,
		},
	}

	for _, tc := range cases {
		mockClient := &mockSGSdkClient{response: tc.expectedByte}
		mockClient.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, nil)
		c := client.NewClient(option.WithHTTPClient(&http.Client{Transport: mockClient}))
		cmd := NewArtifactsCmd(c)
		cmd.SetArgs([]string{
			"list",
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

		var actualResponse api.GeneratedWorkflowListAllArtifactsResponse
		err = json.Unmarshal(out, &actualResponse)
		if err != nil {
			t.Fatal(err)
		}

		if actualResponse.Msg != tc.expectedStruct.Msg {
			t.Fatalf("expected \"%s\" got \"%s\"",
				tc.expectedStruct.Msg,
				actualResponse.Msg)
		}

		if reflect.DeepEqual(actualResponse.Data.Artifacts["url"], tc.expectedStruct.Data.Artifacts["url"]) == false {
			t.Fatalf("expected \"%v\" \ngot \"%v\"",
				tc.expectedStruct.Data.Artifacts["url"],
				actualResponse.Data.Artifacts["url"])
		}

	}
}
