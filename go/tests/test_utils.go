package tests

import (
	"bytes"
	"io"
	"net/http"
	"os/exec"

	"github.com/stretchr/testify/mock"
)

const (
	samplePayloadsDir = "../sample_payloads"
	binaryPath        = "../sg-cli" // Adjust this path based on where your binary is built

	// Command and actions
	cmdArtifacts  = "artifacts"
	cmdWorkflow   = "workflow"
	cmdStack      = "stack"
	actionCreate  = "create"
	actionRead    = "read"
	actionList    = "list"
	actionDelete  = "delete"
	actionApply   = "apply"
	actionDestroy = "destroy"
	actionOutputs = "outputs"

	// Flags
	flagWorkflowGroup = "--workflow-group"
	flagOrg           = "--org"
	flagWorkflowID    = "--workflow-id"
	flagOutputJson    = "--output-json"
	flagStackID       = "--stack-id"
	flagPatchPayload  = "--patch-payload"
	flagForceDelete   = "--force-delete"
	flagBulk          = "--bulk"
	flagRun           = "--run"

	// Error messages
	errOrgNotExist           = "401" // invalid org returns 401
	errRequiredFlag          = "required flag"
	errNoArtifacts           = "No artifacts found for this workflow"
	errStackNotExist         = "Stack does not exist"
	errWfNotExist            = "400" // Stack not found returns 400
	errInvalidJson           = "Error unmarshalling patch JSON"
	errNoSuchFile            = "no such file"
	errMissingPayload        = "Error: accepts 1 arg(s), received 0"
	errInvalidBulkJson       = "json: cannot unmarshal object into Go value of type []create.BulkWorkflow"
	errorMissingResourceName = "Workflow ResourceName is required in object payload, skipping" // No file present
	errStackNotEmpty         = "this stack cannot be deleted since it contains workflows"      // Stack contains workflows

)

type mockSGSdkClient struct {
	mock.Mock
	response   []byte
	statusCode int
}

func (m *mockSGSdkClient) RoundTrip(request *http.Request) (*http.Response, error) {

	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(m.response)),
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}, nil
}

// Helper function to run CLI commands
func runCommand(binaryPath string, args []string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
