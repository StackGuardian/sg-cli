package create

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

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
	// createCmd represents the create command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new workflow",
		Long:  `Create new workflow in the specified organization and workflow group.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Set the options from the command line flags
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.Payload = cmd.Flags().Lookup("payload").Value.String()

			DASHBOARD_URL := "https://app.stackguardian.io/orchestrator"

			payload, err := os.ReadFile(opts.Payload)
			if err != nil {
				cmd.PrintErrln(err)
			}
			if opts.Bulk {
				// Unmarshal the array payload into a slice of BulkWorkflow objects
				var createBulkWorkflowRequest []BulkWorkflow
				err := json.Unmarshal(payload, &createBulkWorkflowRequest)
				if err != nil {
					cmd.Println("Please provide a valid JSON payload. Bulk Payload should be an array of objects.")
					cmd.PrintErrln(err)
					os.Exit(-1)
				}

				// tempMap is needed to delete the CLIConfiguration field and create the workflow request
				var tempMap []map[string]interface{}
				err = json.Unmarshal(payload, &tempMap)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}

				// Iterate over the slice of BulkWorkflow objects
				for idx, bulkWorkflow := range createBulkWorkflowRequest {
					var individualWorkflow *sggosdk.Workflow
					// delete the CLIConfiguration field from the tempMap
					delete(tempMap[idx], "CLIConfiguration")
					// Convert map to []byte and then unmarshal to workflow object
					jsonBody, err := json.Marshal(tempMap[idx])
					if err != nil {
						cmd.PrintErrln(err)
						continue
					}
					err = json.Unmarshal(jsonBody, &individualWorkflow)
					if err != nil {
						cmd.PrintErrln(err)
						continue
					}
					cmd.Println(">> Processing workflow: " + *individualWorkflow.ResourceName)
					err = performPreExecutionFlagChecks(cmd, individualWorkflow, opts)
					if err != nil {
						cmd.PrintErrln(err)
						continue
					}
					// If the workflow group is provided in the bulk payload, use it. Otherwise, use the one provided in the command
					if bulkWorkflow.CLIConfiguration.CLIConfiguration.WorkflowGroup.Name != "" {
						opts.WfgGrp = bulkWorkflow.CLIConfiguration.CLIConfiguration.WorkflowGroup.Name
					} else {
						opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
					}
					response, err := c.Workflows.Create(
						context.Background(),
						opts.Org,
						opts.WfgGrp,
						individualWorkflow,
					)
					if err != nil {
						if !strings.Contains(err.Error(), "Workflow name not unique") {
							cmd.PrintErrln(">> [ERROR] Processing workflow failed for resource name: " + *individualWorkflow.ResourceName + "\n")
							cmd.PrintErrln(err)
							continue
						} else {
							cmd.Println("Workflow already exists, updating instead...")
							// convert to update workflow request
							var updateIndividualWorkflow *sggosdk.PatchedWorkflow
							err = json.Unmarshal(jsonBody, &individualWorkflow)
							if err != nil {
								cmd.PrintErrln(err)
								continue
							}
							response, err := c.Workflows.Patch(
								context.Background(),
								opts.Org,
								*individualWorkflow.ResourceName,
								opts.WfgGrp,
								updateIndividualWorkflow,
							)
							if err != nil {
								cmd.PrintErrln(">> [ERROR] Updating workflow failed for resource name: " + *updateIndividualWorkflow.ResourceName + "\n")
								cmd.PrintErrln(err)
								continue
							}
							if opts.OutputJson {
								cmd.Println(response)
							}
							cmd.Println("Workflow updated successfully.")

							if bulkWorkflow.CLIConfiguration.CLIConfiguration.TfStateFilePath == "" {
								cmd.Println("TfStateFilePath is not provided for workflow: " + *bulkWorkflow.ResourceName)
								cmd.Println(">> Skipping update of state file..\n")
							} else {
								cmd.Println(">> Attempting to upload state file..")
								err = uploadTfState(cmd, &bulkWorkflow, opts)
								if err != nil {
									cmd.PrintErrln("Failed to upload state file for workflow: " + *individualWorkflow.ResourceName + "\n")
									continue
								}
							}
						}
					} else {
						if opts.OutputJson {
							cmd.Println(response)
						}
						cmd.Println("Workflow created successfully.")
						if bulkWorkflow.CLIConfiguration.CLIConfiguration.TfStateFilePath == "" {
							cmd.PrintErrln("[ERROR] TfStateFilePath is not provided for workflow: " + *bulkWorkflow.ResourceName)
							cmd.PrintErrln(">> Skipping update of state file..")
						} else {
							cmd.Println(">> Attempting to upload state file..")
							err = uploadTfState(cmd, &bulkWorkflow, opts)
							if err != nil {
								cmd.PrintErrln("Failed to upload state file for workflow: " + *individualWorkflow.ResourceName + "\n")
							}
						}
						// Run on create
						if opts.Run {
							var createWorkflowRunRequest *sggosdk.WorkflowRun
							err = json.Unmarshal(jsonBody, &createWorkflowRunRequest)
							if err != nil {
								cmd.PrintErrln(err)
								continue
							}
							response, err := c.WorkflowRuns.CreateWorkflowRun(
								context.Background(),
								opts.Org,
								*bulkWorkflow.ResourceName,
								opts.WfgGrp,
								createWorkflowRunRequest,
							)
							if err != nil {
								cmd.PrintErrln("== Failed To Create Workflow Run ==")
								cmd.PrintErrln(err)
								continue
							}
							if opts.OutputJson {
								cmd.Println(response)
							}
							cmd.Println("Workflow run created successfully.")
							workflowRunPath := DASHBOARD_URL +
								"/orgs/" +
								opts.Org +
								"/wfgrps/" +
								opts.WfgGrp +
								"/wfs/" +
								*bulkWorkflow.ResourceName +
								"?tab=runs"
							cmd.Println("To view the workflow run, please visit the following URL:")
							cmd.Println(workflowRunPath)
							//new line for formatting
							cmd.Println()
						}
					}
				}
			} else {
				var createWorkflowRequest *sggosdk.Workflow
				var createWorkflowRunRequest *sggosdk.WorkflowRun
				if opts.PatchPayload != "" && !opts.Bulk {
					patchedJson := utilities.PatchJSON(string(payload), opts.PatchPayload)

					err := json.Unmarshal(
						[]byte(patchedJson),
						&createWorkflowRequest)
					if err != nil {
						cmd.Printf("Error during patching Workflowpayload: %s\n", err)
						os.Exit(-1)
					}
					//unmarshal patched workflow run
					if opts.Run {
						err := json.Unmarshal(
							[]byte(patchedJson),
							&createWorkflowRunRequest)
						if err != nil {
							cmd.Printf("Error during patching WorkflowRun payload: %s\n", err)
							os.Exit(-1)
						}
					}
				} else {
					err := json.Unmarshal(
						payload,
						&createWorkflowRequest)
					if err != nil {
						cmd.Printf("Error while unmarshalling Workflow payload: %s\n", err)
						os.Exit(-1)
					}

					// Ummarshal unpatched workflow run
					if opts.Run {
						err := json.Unmarshal(
							payload,
							&createWorkflowRunRequest)
						if err != nil {
							cmd.Printf("Error while unmarshalling WorkflowRun payload: %s\n", err)
							os.Exit(-1)
						}
					}
				}

				// Perform actions based on the set flags
				if err := performPreExecutionFlagChecks(cmd, createWorkflowRequest, opts); err != nil {
					cmd.PrintErrln(err)
					os.Exit(-1)
				}

				// Run on create
				if opts.Run {
					response, err := c.WorkflowRuns.CreateWorkflowRun(
						context.Background(),
						opts.Org,
						*createWorkflowRequest.ResourceName,
						opts.WfgGrp,
						createWorkflowRunRequest,
					)
					if err != nil {
						cmd.PrintErrln("== Failed To Create Workflow Run ==")
						cmd.PrintErrln(err)
						os.Exit(-1)
					}
					if opts.OutputJson {
						cmd.Println(response)
					}
					cmd.Println("Workflow run created successfully.")
					workflowRunPath := DASHBOARD_URL +
						"/orgs/" +
						opts.Org +
						"/wfgrps/" +
						opts.WfgGrp +
						"/wfs/" +
						*createWorkflowRequest.ResourceName +
						"?tab=runs"
					cmd.Println("To view the workflow run, please visit the following URL:")
					cmd.Println(workflowRunPath)
				} else {
					response, err := c.Workflows.Create(
						context.Background(),
						opts.Org,
						opts.WfgGrp,
						createWorkflowRequest,
					)
					if err != nil {
						cmd.PrintErrln("== Failed To Create Workflow ==")
						cmd.PrintErrln(err)
						os.Exit(-1)
					}
					if opts.OutputJson {
						cmd.Println(response)
					}
					cmd.Println("Workflow created successfully.")
				}
			}

		},
	}

	// Define the flags for the command
	createCmd.Flags().StringVar(&opts.Payload, "payload", "", "The payload JSON file that defines the workflow.")
	createCmd.MarkFlagRequired("payload")

	createCmd.Flags().StringVar(&opts.PatchPayload, "patch-payload", "", "Patch original payload.json input. Add or replace values. Requires valid JSON input.")

	createCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")

	createCmd.Flags().BoolVar(&opts.Preview, "preview", false, "Preview payload content before applying. Execution will not pause.")

	createCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Similar to --preview. But execution will stop, nothing will be applied.")

	createCmd.Flags().BoolVar(&opts.Bulk, "bulk", false, "Bulk import multiple workflows from JSON payload. Upload state files if they exist. Add --run flag to execute")

	createCmd.Flags().BoolVar(&opts.Run, "run", false, "Executes the workflow. Used together with --bulk.")

	return createCmd
}

// performPreExecutionFlagChecks performs pre-execution flag checks and returns the payload
func performPreExecutionFlagChecks(cmd *cobra.Command, payload *sggosdk.Workflow, opts *RunOptions) error {

	if payload.ResourceName == nil || *payload.ResourceName == "" {
		return errors.New(">> [ERROR] Workflow ResourceName is required in object payload, skipping")
	}

	if opts.DryRun {
		requestJson, err := json.MarshalIndent(payload, "", "    ")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(-1)
		}
		cmd.Println(string(requestJson))
		os.Exit(-1)
	} else if opts.Preview {
		requestJson, err := json.MarshalIndent(payload, "", "    ")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(-1)
		}
		cmd.Println(string(requestJson))
	}
	return nil
}

// uploadTfState uploads the Terraform state file to Stackguardian
func uploadTfState(cmd *cobra.Command, payload *BulkWorkflow, opts *RunOptions) error {
	SG_BASE_URL := os.Getenv("SG_BASE_URL")
	SG_API_TOKEN := os.Getenv("SG_API_TOKEN")

	// Get the tfstate upload url for the workflow
	url := SG_BASE_URL + "/api/v1/orgs/" +
		opts.Org +
		"/wfgrps/" +
		opts.WfgGrp +
		"/wfs/" +
		*payload.ResourceName +
		"/tfstate_upload_url"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		cmd.PrintErrln(">> [ERROR] Failed to get tfstate upload url for workflow: " + *payload.ResourceName + "\n")
		cmd.PrintErrln(err)
		return err
	}
	req.Header.Set("Authorization", "apikey "+SG_API_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		cmd.PrintErrln(">> [ERROR] Failed to get tfstate upload url for workflow: " + *payload.ResourceName + "\n")
		cmd.PrintErrln(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		cmd.PrintErrln(">> [ERROR] Failed to get tfstate upload url for workflow: " + *payload.ResourceName + "\n")
		cmd.PrintErrln("Expected status code 200, got " + resp.Status)
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cmd.PrintErrln(">> [ERROR] Failed to get tfstate upload url for workflow: " + *payload.ResourceName + "\n")
		cmd.PrintErrln(err)
		return err
	}
	var response tfStateUploadUrlResponse
	if err := json.Unmarshal(body, &response); err != nil {
		cmd.PrintErrln(">> [ERROR] Failed to get tfstate upload url for workflow: " + *payload.ResourceName + "\n")
		cmd.PrintErrln(err)
		return err
	}
	tfUploadUrl := response.Msg

	// Use the tfUploadUrl to upload the state file to Stackguardian
	cmd.Println(">> Uploading state file to Stackguardian..")
	// Create a temporary directory to store the state file
	tmpDir, err := os.MkdirTemp("", opts.Org+"-"+opts.WfgGrp)
	if err != nil {
		cmd.PrintErrln(">> [ERROR] Failed to create temp directory for state file upload: " + *payload.ResourceName + "\n")
		cmd.PrintErrln(err)
		return err
	}
	defer os.RemoveAll(tmpDir)
	copyFileCmd := exec.Command("cp", payload.CLIConfiguration.CLIConfiguration.TfStateFilePath, tmpDir+"/tfstate.json")
	err = copyFileCmd.Run()
	if err != nil {
		cmd.PrintErrln(">> [ERROR] Failed to access state file : " + payload.CLIConfiguration.CLIConfiguration.TfStateFilePath +
			". Please check if the state file exists and is accessible.")
		cmd.PrintErrln(err)
		return err
	}

	// Use curl to upload the state file to the tfUploadUrl
	// TODO: Implement this using Go's native HTTP client
	curlCmd := exec.Command("curl", "-i", "-s", "-X", "PUT",
		"-H", "Accept: application/json, text/plain, */*",
		"-H", "Content-Type: application/json",
		"-H", "ContentType: application/json",
		"-T", tmpDir+"/tfstate.json",
		tfUploadUrl,
	)

	output, err := curlCmd.CombinedOutput()
	if err != nil {
		cmd.PrintErrln(">> [ERROR] Error running curl command:", err)
	}

	if strings.Contains(string(output), "HTTP/1.1 200 OK") {
		cmd.Println(">> State file uploaded successfully.")
	} else {
		cmd.PrintErrln(">> [ERROR] Failed to upload state file for workflow: " + *payload.ResourceName + "\n")
		cmd.PrintErrln(string(output))
	}

	return nil
}
