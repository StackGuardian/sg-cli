package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArtifactsE2E(t *testing.T) {
	// Setup test
	orgName := "demo-org"
	workflowGroup := "test-terragrunt"
	workflowID := "CUSTOM-7OeX" // Replace with a valid workflow ID that has artifacts

	t.Run("Artifacts_List_Basic", func(t *testing.T) {
		// List artifacts for a workflow
		listArgs := []string{
			cmdArtifacts, actionList,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}

		output, err := runCommand(binaryPath, listArgs)
		t.Logf("List artifacts output: %s", output)

		// The test might pass or fail depending on whether the workflow has artifacts
		// We're just checking that the command runs without errors
		if err != nil {
			// If there's an error, check if it's because there are no artifacts
			if assert.Contains(t, output, errNoArtifacts) {
				t.Log("No artifacts found for this workflow, which is acceptable")
			} else {
				t.Errorf("Expected no error or 'No artifacts' message, got: %v", err)
			}
		} else {
			// If no error, the command succeeded and we should have some output
			assert.Contains(t, output, "artifacts")
		}
	})

	t.Run("Artifacts_List_With_OutputJson", func(t *testing.T) {
		// List artifacts with output-json flag
		listArgs := []string{
			cmdArtifacts, actionList,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
			flagOutputJson,
		}

		output, err := runCommand(binaryPath, listArgs)
		t.Logf("List artifacts with output-json output: %s", output)

		// The test might pass or fail depending on whether the workflow has artifacts
		if err != nil {
			// If there's an error, check if it's because there are no artifacts
			if assert.Contains(t, output, errNoArtifacts) {
				t.Log("No artifacts found for this workflow, which is acceptable")
			} else {
				t.Errorf("Expected no error or 'No artifacts' message, got: %v", err)
			}
		} else {
			// If no error, the command succeeded and we should have JSON output
			assert.Contains(t, output, "{")
			assert.Contains(t, output, "}")
		}
	})

	t.Run("Negative_Tests-Invalid_Organization", func(t *testing.T) {
		invalidOrg := "non-existent-org"
		// Expected to return 401 Unauthorized
		listArgs := []string{
			cmdArtifacts, actionList,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err := runCommand(binaryPath, listArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)
	})

	t.Run("Negative_Tests-Invalid_WorkflowGroup", func(t *testing.T) {
		invalidWfGroup := "non-existent-workflow-group"
		listArgs := []string{
			cmdArtifacts, actionList,
			flagOrg, orgName,
			flagWorkflowGroup, invalidWfGroup,
			flagWorkflowID, workflowID,
		}
		output, err := runCommand(binaryPath, listArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errNoArtifacts)
	})

	t.Run("Negative_Tests-Invalid_WorkflowID", func(t *testing.T) {
		invalidWfID := "non-existent-workflow"
		listArgs := []string{
			cmdArtifacts, actionList,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, invalidWfID,
		}
		output, err := runCommand(binaryPath, listArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errNoArtifacts)
	})

	t.Run("Negative_Tests-Missing_Required_Flags", func(t *testing.T) {
		// List without org
		listArgs := []string{
			cmdArtifacts, actionList,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err := runCommand(binaryPath, listArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)

		// List without workflow group
		listArgs = []string{
			cmdArtifacts, actionList,
			flagOrg, orgName,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, listArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)

		// List without workflow ID
		listArgs = []string{
			cmdArtifacts, actionList,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, listArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)
	})
}
