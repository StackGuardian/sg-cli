package create

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/StackGuardian/sg-cli/cmd/output"
	"github.com/StackGuardian/sg-cli/utilities"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

// Extend the Workflow struct from the sdk to add the new field for bulk
type BulkWorkflow struct {
	sggosdk.Workflow
	CLIConfiguration
}

// Without an additional wrapper json.Unmarshal will not unmarshal the nested structs
type CLIConfiguration struct {
	CLIConfiguration CLIConfigurationStruct `json:"CLIConfiguration"`
}

type CLIConfigurationStruct struct {
	WorkflowGroup   WorkflowGroup `json:"WorkflowGroup"`
	TfStateFilePath string        `json:"TfStateFilePath"`
}

type WorkflowGroup struct {
	Name string `json:"name"`
}

type RunOptions struct {
	Org          string
	WfgGrp       string
	Bulk         bool
	Preview      bool
	DryRun       bool
	Run          bool
	OutputJson   bool
	PatchPayload string
	Payload      string
}

type tfStateUploadUrlResponse struct {
	Msg string `json:"msg"`
}

func (o *BulkWorkflow) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &o.Workflow); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &o.CLIConfiguration); err != nil {
		return err
	}
	return nil
}

func NewCreateCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	const dashboardURL = "https://app.stackguardian.io/orchestrator"

	var createCmd = &cobra.Command{
		Use:   "create [payload.json]",
		Short: "Create a new workflow",
		Long:  "Create a new workflow from a JSON payload file. Supports bulk mode, dry-run preview, and patching.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.Payload = args[0]

			payload, err := os.ReadFile(opts.Payload)
			if err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}

			if opts.Bulk {
				runBulkCreate(cmd, c, payload, opts, dashboardURL)
			} else {
				runSingleCreate(cmd, c, payload, opts, dashboardURL)
			}
		},
	}

	createCmd.Flags().StringVar(&opts.PatchPayload, "patch-payload", "", "Merge a JSON patch over the payload before applying.")
	createCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")
	createCmd.Flags().BoolVar(&opts.Preview, "preview", false, "Preview the payload before applying (execution continues).")
	createCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the payload and exit without applying.")
	createCmd.Flags().BoolVar(&opts.Bulk, "bulk", false, "Bulk-create workflows from a JSON array payload.")
	createCmd.Flags().BoolVar(&opts.Run, "run", false, "Trigger a run immediately after creation.")

	return createCmd
}

func cmdPrintln(cmd *cobra.Command, v interface{}) {
	if cmd != nil {
		cmd.Println(v)
	} else {
		fmt.Println(v)
	}
}

// RunCreate is the public entry point for interactive/programmatic use.
func RunCreate(c *client.Client, opts *RunOptions) {
	const url = "https://app.stackguardian.io/orchestrator"
	payload, err := os.ReadFile(opts.Payload)
	if err != nil {
		output.Error(err.Error())
		return
	}
	if opts.Bulk {
		runBulkCreate(nil, c, payload, opts, url)
	} else {
		runSingleCreate(nil, c, payload, opts, url)
	}
}

func runBulkCreate(cmd *cobra.Command, c *client.Client, payload []byte, opts *RunOptions, dashboardURL string) {
	var createBulkWorkflowRequest []BulkWorkflow
	if err := json.Unmarshal(payload, &createBulkWorkflowRequest); err != nil {
		output.Error("Invalid JSON payload — bulk payload must be an array of objects.")
		output.Error(err.Error())
		os.Exit(1)
	}

	var tempMap []map[string]interface{}
	if err := json.Unmarshal(payload, &tempMap); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}

	for idx, bulkWorkflow := range createBulkWorkflowRequest {
		var individualWorkflow *sggosdk.Workflow
		delete(tempMap[idx], "CLIConfiguration")
		jsonBody, err := json.Marshal(tempMap[idx])
		if err != nil {
			output.Error(err.Error())
			continue
		}
		if err := json.Unmarshal(jsonBody, &individualWorkflow); err != nil {
			output.Error(err.Error())
			continue
		}

		output.Info("Processing workflow: " + individualWorkflow.ResourceName.Value)

		if err := performPreExecutionFlagChecks(cmd, individualWorkflow, opts); err != nil {
			output.Error(err.Error())
			continue
		}

		wfGrp := opts.WfgGrp
		if bulkWorkflow.CLIConfiguration.CLIConfiguration.WorkflowGroup.Name != "" {
			wfGrp = bulkWorkflow.CLIConfiguration.CLIConfiguration.WorkflowGroup.Name
		}

		var response interface{}
		apiErr := output.WithSpinner("Creating "+individualWorkflow.ResourceName.Value+"...", func() error {
			var e error
			response, e = c.Workflows.CreateWorkflow(context.Background(), opts.Org, wfGrp, individualWorkflow)
			return e
		})

		if apiErr != nil {
			if !strings.Contains(apiErr.Error(), "Workflow name not unique") {
				output.Error("Failed to create " + individualWorkflow.ResourceName.Value + ": " + apiErr.Error())
				continue
			}
			output.Warning("Workflow already exists — updating instead...")
			var updateIndividualWorkflow *sggosdk.PatchedWorkflow
			if err := json.Unmarshal(jsonBody, &individualWorkflow); err != nil {
				output.Error(err.Error())
				continue
			}
			var updateResponse interface{}
			updateErr := output.WithSpinner("Updating "+individualWorkflow.ResourceName.Value+"...", func() error {
				var e error
				updateResponse, e = c.Workflows.UpdateWorkflow(context.Background(), opts.Org, individualWorkflow.ResourceName.Value, wfGrp, updateIndividualWorkflow)
				return e
			})
			if updateErr != nil {
				output.Error("Failed to update " + *updateIndividualWorkflow.ResourceName + ": " + updateErr.Error())
				continue
			}
			if opts.OutputJson {
				cmdPrintln(cmd, updateResponse)
			}
			output.Success("Workflow updated successfully.")
			handleTfState(cmd, &bulkWorkflow, opts)
			continue
		}

		if opts.OutputJson {
			cmdPrintln(cmd, response)
		}
		output.Success("Workflow created successfully.")
		handleTfState(cmd, &bulkWorkflow, opts)

		if opts.Run {
			runWorkflow(cmd, c, &bulkWorkflow, jsonBody, opts, dashboardURL, wfGrp)
		}
	}
}

func handleTfState(cmd *cobra.Command, bulkWorkflow *BulkWorkflow, opts *RunOptions) {
	if bulkWorkflow.CLIConfiguration.CLIConfiguration.TfStateFilePath == "" {
		output.Warning("TfStateFilePath not provided for " + bulkWorkflow.ResourceName.Value + " — skipping state upload.")
		return
	}
	output.Info("Uploading Terraform state file...")
	if err := uploadTfState(cmd, bulkWorkflow, opts); err != nil {
		output.Error("Failed to upload state file for " + bulkWorkflow.ResourceName.Value + ": " + err.Error())
	}
}

func runWorkflow(cmd *cobra.Command, c *client.Client, bulkWorkflow *BulkWorkflow, jsonBody []byte, opts *RunOptions, dashboardURL, wfGrp string) {
	var createWorkflowRunRequest *sggosdk.WorkflowRun
	if err := json.Unmarshal(jsonBody, &createWorkflowRunRequest); err != nil {
		output.Error(err.Error())
		return
	}
	var runResponse interface{}
	err := output.WithSpinner("Triggering workflow run...", func() error {
		var e error
		runResponse, e = c.WorkflowRuns.CreateWorkflowRun(context.Background(), opts.Org, bulkWorkflow.ResourceName.Value, wfGrp, createWorkflowRunRequest)
		return e
	})
	if err != nil {
		output.Error("Failed to create workflow run: " + err.Error())
		return
	}
	if opts.OutputJson {
		cmdPrintln(cmd, runResponse)
	}
	output.Success("Workflow run created successfully.")
	output.URL("View run at:", dashboardURL+"/orgs/"+opts.Org+"/wfgrps/"+wfGrp+"/wfs/"+bulkWorkflow.ResourceName.Value+"?tab=runs")
	output.Newline()
}

func runSingleCreate(cmd *cobra.Command, c *client.Client, payload []byte, opts *RunOptions, dashboardURL string) {
	var createWorkflowRequest *sggosdk.Workflow
	var createWorkflowRunRequest *sggosdk.WorkflowRun

	if opts.PatchPayload != "" {
		patched := utilities.PatchJSON(string(payload), opts.PatchPayload)
		if err := json.Unmarshal([]byte(patched), &createWorkflowRequest); err != nil {
			output.Error("Error patching workflow payload: " + err.Error())
			os.Exit(1)
		}
		if opts.Run {
			if err := json.Unmarshal([]byte(patched), &createWorkflowRunRequest); err != nil {
				output.Error("Error patching workflow run payload: " + err.Error())
				os.Exit(1)
			}
		}
	} else {
		if err := json.Unmarshal(payload, &createWorkflowRequest); err != nil {
			output.Error("Error reading workflow payload: " + err.Error())
			os.Exit(1)
		}
		if opts.Run {
			if err := json.Unmarshal(payload, &createWorkflowRunRequest); err != nil {
				output.Error("Error reading workflow run payload: " + err.Error())
				os.Exit(1)
			}
		}
	}

	if err := performPreExecutionFlagChecks(cmd, createWorkflowRequest, opts); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}

	if opts.Run {
		var response interface{}
		err := output.WithSpinner("Creating workflow run...", func() error {
			var e error
			response, e = c.WorkflowRuns.CreateWorkflowRun(context.Background(), opts.Org, createWorkflowRequest.ResourceName.Value, opts.WfgGrp, createWorkflowRunRequest)
			return e
		})
		if err != nil {
			output.Error("Failed to create workflow run: " + err.Error())
			os.Exit(1)
		}
		if opts.OutputJson {
			cmdPrintln(cmd, response)
		}
		output.Success("Workflow run created successfully.")
		output.URL("View run at:", dashboardURL+"/orgs/"+opts.Org+"/wfgrps/"+opts.WfgGrp+"/wfs/"+createWorkflowRequest.ResourceName.Value+"?tab=runs")
	} else {
		var response interface{}
		err := output.WithSpinner("Creating workflow...", func() error {
			var e error
			response, e = c.Workflows.CreateWorkflow(context.Background(), opts.Org, opts.WfgGrp, createWorkflowRequest)
			return e
		})
		if err != nil {
			output.Error("Failed to create workflow: " + err.Error())
			os.Exit(1)
		}
		if opts.OutputJson {
			cmdPrintln(cmd, response)
		}
		output.Success("Workflow created successfully.")
	}
}

// performPreExecutionFlagChecks validates and optionally previews the payload.
func performPreExecutionFlagChecks(cmd *cobra.Command, payload *sggosdk.Workflow, opts *RunOptions) error {
	if payload.ResourceName == nil || payload.ResourceName.Value == "" {
		return errors.New("workflow ResourceName is required in the payload")
	}

	if opts.DryRun || opts.Preview {
		requestJson, err := json.MarshalIndent(payload, "", "    ")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}
		output.Section("Payload Preview")
		cmdPrintln(cmd, string(requestJson))
		if opts.DryRun {
			os.Exit(0)
		}
	}
	return nil
}

// uploadTfState uploads the Terraform state file to StackGuardian.
func uploadTfState(cmd *cobra.Command, payload *BulkWorkflow, opts *RunOptions) error {
	SG_BASE_URL := os.Getenv("SG_BASE_URL")
	SG_API_TOKEN := os.Getenv("SG_API_TOKEN")

	url := SG_BASE_URL + "/api/v1/orgs/" + opts.Org + "/wfgrps/" + opts.WfgGrp + "/wfs/" + payload.ResourceName.Value + "/tfstate_upload_url"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "apikey "+SG_API_TOKEN)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("expected HTTP 200, got " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var uploadResp tfStateUploadUrlResponse
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return err
	}
	tfUploadUrl := uploadResp.Msg

	output.Info("Uploading state file to StackGuardian...")

	tmpDir, err := os.MkdirTemp("", opts.Org+"-"+opts.WfgGrp)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if err := exec.Command("cp", payload.CLIConfiguration.CLIConfiguration.TfStateFilePath, tmpDir+"/tfstate.json").Run(); err != nil {
		return errors.New("cannot access state file at " + payload.CLIConfiguration.CLIConfiguration.TfStateFilePath)
	}

	// TODO: replace curl with native Go HTTP PUT
	curlCmd := exec.Command("curl", "-i", "-s", "-X", "PUT",
		"-H", "Accept: application/json, text/plain, */*",
		"-H", "Content-Type: application/json",
		"-H", "ContentType: application/json",
		"-T", tmpDir+"/tfstate.json",
		tfUploadUrl,
	)

	uploadOut, err := curlCmd.CombinedOutput()
	if err != nil {
		return err
	}

	if strings.Contains(string(uploadOut), "HTTP/1.1 200 OK") {
		output.Success("State file uploaded successfully.")
	} else {
		return errors.New("state file upload failed: " + string(uploadOut))
	}

	return nil
}
