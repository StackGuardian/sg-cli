package tests

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	createWorkflowFile              = "create_workflow.json"
	missingResourceNameWorkflowFile = "missing_resource_name.json"
	invalidJsonWorkflowFile         = "invalid_json_workflow.json"
	bulkWorkflowFile                = "create_bulk_workflow.json"
	bulkWorkflowIdPrefix            = "bulk-workflow-create-test"

	// Informational messages
	msgWorkflowAlreadyExist           = "Workflow already exists, updating instead"
	msgWorkflowCreatedSuccessfully    = "Workflow created successfully."
	msgWorkflowUpdatedSuccessfully    = "Workflow updated successfully."
	msgWorkflowRunCreatedSuccessfully = "Workflow run created successfully."
	msgWorkflowAppliedSuccessfully    = "Workflow apply run successfully."
	msgWorkflowDestroyedSuccessfully  = "Workflow destroy run successfully."
	msgWorkflowDeletedSuccessfully    = "Workflow deleted successfully."
)

func generateWorkflowName() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("test-workflow-%d-%d", time.Now().Unix(), n.Int64())
}

func TestWorkflowE2E(t *testing.T) {
	// Setup test
	orgName := "demo-org"
	workflowGroup := "sg-sdk-go-test"
	workflowID := generateWorkflowName()

	t.Run("Workflow_Basic_Operations", func(t *testing.T) {
		// Step 1: Create workflow
		samplePayloadPath := filepath.Join(samplePayloadsDir, createWorkflowFile)
		patchPayload := fmt.Sprintf(`{"ResourceName":"%s"}`, workflowID)

		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, patchPayload,
			"--", samplePayloadPath,
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgWorkflowCreatedSuccessfully)

		// Step 2: Read workflow
		readArgs := []string{
			cmdWorkflow, actionRead,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, readArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, workflowID)

		// Step 3: List workflows
		listArgs := []string{
			cmdWorkflow, actionList,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, listArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, workflowID)

		// Step 4: Apply workflow
		applyArgs := []string{
			cmdWorkflow, actionApply,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, applyArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgWorkflowAppliedSuccessfully)

		// Step 5: Destroy workflow
		destroyArgs := []string{
			cmdWorkflow, actionDestroy,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, destroyArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgWorkflowDestroyedSuccessfully)

		// Step 6: Delete workflow (cleanup)
		deleteArgs := []string{
			cmdWorkflow, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgWorkflowDeletedSuccessfully)

		// Verify workflow is deleted
		output, err = runCommand(binaryPath, readArgs)
		assert.Error(t, err, "Expected error when reading deleted workflow")
		t.Logf("Read after delete output: %s", output)
	})

	t.Run("Negative_Tests-Invalid_Organization", func(t *testing.T) {
		invalidOrg := "non-existent-org"
		// Expected to return 401 Unauthorized

		// Create with invalid org
		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			"--", filepath.Join(samplePayloadsDir, createWorkflowFile),
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// Read with invalid org
		readArgs := []string{
			cmdWorkflow, actionRead,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, readArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// Apply with invalid org
		applyArgs := []string{
			cmdWorkflow, actionApply,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, applyArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// Destroy with invalid org
		destroyArgs := []string{
			cmdWorkflow, actionDestroy,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, workflowID,
		}
		output, err = runCommand(binaryPath, destroyArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// List with invalid org
		listArgs := []string{
			cmdWorkflow, actionList,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, listArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)
	})

	t.Run("Negative_Tests-Read_Invalid_Workflow_ID", func(t *testing.T) {
		invalidID := "non-existent-workflow-id"

		// Read non-existent workflow
		readArgs := []string{
			cmdWorkflow, actionRead,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, invalidID,
		}
		output, err := runCommand(binaryPath, readArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errWfNotExist)
	})

	t.Run("Negative_Tests-Missing_Required_Flags", func(t *testing.T) {
		// Create without org
		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagWorkflowGroup, workflowGroup,
			"--", filepath.Join(samplePayloadsDir, createWorkflowFile),
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)

		// Create without workflow group
		createArgs = []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			"--", filepath.Join(samplePayloadsDir, createWorkflowFile),
		}
		output, err = runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)

		// Create without payload file
		createArgs = []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errMissingPayload)

		// Apply without workflow ID
		applyArgs := []string{
			cmdWorkflow, actionApply,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, applyArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)
	})

	t.Run("Negative_Tests-Invalid_Patch_Payload", func(t *testing.T) {
		// Create with invalid JSON in patch payload
		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, "{invalid-json}",
			"--", filepath.Join(samplePayloadsDir, createWorkflowFile),
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errInvalidJson)
	})

	t.Run("Negative_Tests-Non_Existent_Payload_File", func(t *testing.T) {
		// Create with non-existent file
		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			"--", "non-existent-file.json",
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errNoSuchFile)
	})

	t.Run("Negative_Tests-Missing_ResourceName_In_Payload", func(t *testing.T) {
		// Create with invalid JSON payload
		// This create request will fail with the error "Workflow ResourceName is required in object payload, skipping"
		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			"--", filepath.Join(samplePayloadsDir, missingResourceNameWorkflowFile),
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errorMissingResourceName)
	})

	t.Run("Workflow_Bulk_Create", func(t *testing.T) {
		// Create workflows in bulk
		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagBulk,
			flagRun,
			"--", filepath.Join(samplePayloadsDir, bulkWorkflowFile),
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		// Check the bulk creation outputs
		// The output should not contain messages that the workflows already exist or any updations
		assert.NotContains(t, output, msgWorkflowAlreadyExist)
		assert.NotContains(t, output, msgWorkflowUpdatedSuccessfully)
		// Check to make sure that there are exactly 3 new workflows created.
		workflowCreatedCount := strings.Count(output, msgWorkflowCreatedSuccessfully)
		assert.Equal(t, 3, workflowCreatedCount)
		//Check to make sure the all 3 workflows have run due to the --run flag
		workflowRunCount := strings.Count(output, msgWorkflowRunCreatedSuccessfully)
		assert.Equal(t, 3, workflowRunCount)
		// Check for messages that are expected based on the bulk payload
		assert.Contains(t, output, errorMissingResourceName)
		assert.Contains(t, output, "Failed to access state file : ../sample_payloads/nonExistent.tfstate. "+
			"Please check if the state file exists and is accessible.")
		assert.Contains(t, output, "Failed to upload state file for workflow: "+bulkWorkflowIdPrefix+"-1")
		assert.Contains(t, output, "TfStateFilePath is not provided for workflow: "+bulkWorkflowIdPrefix+"-2")
		assert.Contains(t, output, "Skipping update of state file..")
		assert.Contains(t, output, "State file uploaded successfully.")

		// verify creation using read
		readArgs := []string{
			cmdWorkflow, actionRead,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-1",
		}
		output, err = runCommand(binaryPath, readArgs)
		assert.NoError(t, err)
		t.Logf("Read output: %s", output)
		assert.Contains(t, output, bulkWorkflowIdPrefix+"-1")

		readArgs = []string{
			cmdWorkflow, actionRead,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-2",
		}
		output, err = runCommand(binaryPath, readArgs)
		assert.NoError(t, err)
		t.Logf("Read output: %s", output)
		assert.Contains(t, output, bulkWorkflowIdPrefix+"-2")

		readArgs = []string{
			cmdWorkflow, actionRead,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-3",
		}
		output, err = runCommand(binaryPath, readArgs)
		assert.NoError(t, err)
		t.Logf("Read output: %s", output)
		assert.Contains(t, output, bulkWorkflowIdPrefix+"-3")

		// Lets run the create again, this time it find that the workflows already exist
		// and try to update them
		createArgs = []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagBulk,
			"--", filepath.Join(samplePayloadsDir, bulkWorkflowFile),
		}

		output, err = runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		// Check the bulk creation outputs
		// No new workflows should be created
		assert.NotContains(t, output, msgWorkflowCreatedSuccessfully)
		// There should be 3 updated workflows
		workflowUpdatedCount := strings.Count(output, msgWorkflowUpdatedSuccessfully)
		assert.Equal(t, 3, workflowUpdatedCount)
		assert.Contains(t, output, msgWorkflowAlreadyExist)
		// Checks to make sure the bulk request is as expected.
		assert.Contains(t, output, errorMissingResourceName)
		assert.Contains(t, output, "Failed to access state file : ../sample_payloads/nonExistent.tfstate. "+
			"Please check if the state file exists and is accessible.")
		assert.Contains(t, output, "Failed to upload state file for workflow: "+bulkWorkflowIdPrefix+"-1")
		assert.Contains(t, output, "TfStateFilePath is not provided for workflow: "+bulkWorkflowIdPrefix+"-2")
		assert.Contains(t, output, "Skipping update of state file..")
		assert.Contains(t, output, "State file uploaded successfully.")

		// Perform clean up using delete
		deleteArgs := []string{
			cmdWorkflow, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-1",
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.NoError(t, err)
		t.Logf("Delete workflow 1 output: %s", output)

		deleteArgs = []string{
			cmdWorkflow, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-2",
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.NoError(t, err)
		t.Logf("Delete workflow 2 output: %s", output)

		deleteArgs = []string{
			cmdWorkflow, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-3",
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.NoError(t, err)
		t.Logf("Delete workflow 3 output: %s", output)

		// Verify deletion using read
		readArgs = []string{
			cmdWorkflow, actionRead,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-1",
		}
		output, err = runCommand(binaryPath, readArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errWfNotExist)

		readArgs = []string{
			cmdWorkflow, actionRead,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagWorkflowID, bulkWorkflowIdPrefix + "-2",
		}
		output, err = runCommand(binaryPath, readArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errWfNotExist)
	})

	t.Run("Negative_Tests-Invalid_Bulk_Payload", func(t *testing.T) {
		// Create workflows in bulk with invalid JSON in patch"
		// Create workflows in bulk
		createArgs := []string{
			cmdWorkflow, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagBulk,
			"--", filepath.Join(samplePayloadsDir, invalidJsonWorkflowFile),
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errInvalidBulkJson)
		t.Logf(output)
	})
}
