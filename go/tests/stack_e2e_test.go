package tests

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	createStackFile                       = "create_stack_qa.json"
	createStackWithWorkflowReferencesFile = "create_stack_with_workflow_refs.json"
	createStackWithWorkflowConfigFile     = "create_stack_workflows_config.json"

	// Informational messages
	msgStackCreatedSuccessfully        = "Stack created successfully."
	msgStackDeletedSuccessfully        = "Stack deleted successfully."
	msgStackAppliedSuccessfully        = "Stack apply executed."
	msgStackDestroyedSuccessfully      = "Stack Workflow destroy run successfully."
	msgStackOutputsFetchedSuccessfully = "Stack output fetched successfully"
)

func generateStackName() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("test-stack-%d-%d", time.Now().Unix(), n.Int64())
}

func TestStackE2E(t *testing.T) {
	// Test setup
	orgName := "demo-org"
	workflowGroup := "sg-sdk-go-test"

	t.Run("Stack_Basic_Operations", func(t *testing.T) {
		// Step 1: Create stack
		samplePayloadPath := filepath.Join(samplePayloadsDir, createStackFile)
		randomStackName := generateStackName()
		patchPayload := fmt.Sprintf(`{"ResourceName":"%s"}`, randomStackName)

		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, patchPayload,
			"--", samplePayloadPath,
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackCreatedSuccessfully)

		// Step 2: Apply stack
		applyArgs := []string{
			cmdStack, actionApply,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
		}
		output, err = runCommand(binaryPath, applyArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackAppliedSuccessfully)

		// Step 3: Get stack outputs
		outputsArgs := []string{
			cmdStack, actionOutputs,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
		}
		output, err = runCommand(binaryPath, outputsArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackOutputsFetchedSuccessfully)

		// Step 4: Destroy stack
		destroyArgs := []string{
			cmdStack, actionDestroy,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
		}
		output, err = runCommand(binaryPath, destroyArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackDestroyedSuccessfully)

		// Step 5: Delete stack (cleanup)
		deleteArgs := []string{
			cmdStack, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
			flagForceDelete,
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackDeletedSuccessfully)
	})

	t.Run("Stack_With_Workflows_Reference", func(t *testing.T) {
		// Create a stack with workflow references
		stackWithWorkflowsID := generateStackName()
		samplePayloadPath := filepath.Join(samplePayloadsDir, createStackWithWorkflowReferencesFile)
		patchPayload := fmt.Sprintf(`{"ResourceName":"%s"}`, stackWithWorkflowsID)

		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, patchPayload,
			"--", samplePayloadPath,
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackCreatedSuccessfully)

		// Delete with force-delete
		deleteWithForceArgs := []string{
			cmdStack, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, stackWithWorkflowsID,
			flagForceDelete,
		}
		output, err = runCommand(binaryPath, deleteWithForceArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackDeletedSuccessfully)
	})

	t.Run("Stack_With_Workflows_Config", func(t *testing.T) {
		// Create a stack with workflow references
		stackWithWorkflowsID := generateStackName()
		samplePayloadPath := filepath.Join(samplePayloadsDir, createStackWithWorkflowConfigFile)
		patchPayload := fmt.Sprintf(`{"ResourceName":"%s"}`, stackWithWorkflowsID)

		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, patchPayload,
			"--", samplePayloadPath,
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackCreatedSuccessfully)

		// Delete with force-delete
		deleteWithForceArgs := []string{
			cmdStack, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, stackWithWorkflowsID,
			flagForceDelete,
		}
		output, err = runCommand(binaryPath, deleteWithForceArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackDeletedSuccessfully)
	})

	t.Run("Output_JSON_Flag", func(t *testing.T) {
		// Create a stack with output-json flag
		outputJsonStackID := generateStackName()
		samplePayloadPath := filepath.Join(samplePayloadsDir, createStackFile)
		patchPayload := fmt.Sprintf(`{"ResourceName":"%s"}`, outputJsonStackID)

		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, patchPayload,
			flagOutputJson,
			"--", samplePayloadPath,
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		// The output should contain JSON
		assert.Contains(t, output, "{")
		assert.Contains(t, output, fmt.Sprintf("Stack %s created", outputJsonStackID))
		assert.Contains(t, output, "}")

		// Clean up - delete the stack
		deleteArgs := []string{
			cmdStack, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, outputJsonStackID,
			flagForceDelete,
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackDeletedSuccessfully)
	})

	t.Run("Test_Force_Delete_Flag", func(t *testing.T) {
		outputJsonStackID := generateStackName()
		samplePayloadPath := filepath.Join(samplePayloadsDir, createStackFile)
		patchPayload := fmt.Sprintf(`{"ResourceName":"%s"}`, outputJsonStackID)

		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, patchPayload,
			"--", samplePayloadPath,
		}

		output, err := runCommand(binaryPath, createArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackCreatedSuccessfully)

		// Try to delete without force-delete (should fail)
		deleteArgs := []string{
			cmdStack, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagStackID, outputJsonStackID,
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errStackNotEmpty)

		// Delete the stack with force-delete
		deleteArgs = append(deleteArgs, flagForceDelete)
		output, err = runCommand(binaryPath, deleteArgs)
		assert.NoError(t, err)
		assert.Contains(t, output, msgStackDeletedSuccessfully)
	})

	t.Run("Negative_Tests-Invalid_Organization", func(t *testing.T) {
		invalidOrg := "non-existent-org"
		// Expected to return 401 Unauthorized
		randomStackName := generateStackName()
		patchPayload := fmt.Sprintf(`{"ResourceName":"%s"}`, randomStackName)
		// Create with invalid org
		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, patchPayload,
			"--", filepath.Join(samplePayloadsDir, createStackFile),
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// Apply with invalid org
		applyArgs := []string{
			cmdStack, actionApply,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
		}
		output, err = runCommand(binaryPath, applyArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// Outputs with invalid org
		outputsArgs := []string{
			cmdStack, actionOutputs,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
		}
		output, err = runCommand(binaryPath, outputsArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// Destroy with invalid org
		destroyArgs := []string{
			cmdStack, actionDestroy,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
		}
		output, err = runCommand(binaryPath, destroyArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)

		// Delete with invalid org
		deleteArgs := []string{
			cmdStack, actionDelete,
			flagOrg, invalidOrg,
			flagWorkflowGroup, workflowGroup,
			flagStackID, randomStackName,
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errOrgNotExist)
	})

	t.Run("Negative_Tests-Missing_Required_Flags", func(t *testing.T) {
		// Create without org
		createArgs := []string{
			cmdStack, actionCreate,
			flagWorkflowGroup, workflowGroup,
			"--", filepath.Join(samplePayloadsDir, createStackFile),
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)

		// Create without workflow group
		createArgs = []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			"--", filepath.Join(samplePayloadsDir, createStackFile),
		}
		output, err = runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)

		// Create without payload file
		createArgs = []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errMissingPayload)

		// Apply without stack ID
		applyArgs := []string{
			cmdStack, actionApply,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, applyArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)

		// Delete without stack ID
		deleteArgs := []string{
			cmdStack, actionDelete,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
		}
		output, err = runCommand(binaryPath, deleteArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errRequiredFlag)
	})

	t.Run("Negative_Tests-Invalid_Patch_Payload", func(t *testing.T) {
		// Create with invalid JSON in patch payload
		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			flagPatchPayload, "{invalid-json}",
			"--", filepath.Join(samplePayloadsDir, createStackFile),
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errInvalidJson)
	})

	t.Run("Negative_Tests-Non_Existent_Payload_File", func(t *testing.T) {
		// Create with non-existent file
		createArgs := []string{
			cmdStack, actionCreate,
			flagOrg, orgName,
			flagWorkflowGroup, workflowGroup,
			"--", "non-existent-file.json",
		}
		output, err := runCommand(binaryPath, createArgs)
		assert.Error(t, err)
		assert.Contains(t, output, errNoSuchFile)
	})
}
