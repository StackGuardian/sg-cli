package create

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/StackGuardian/sg-cli/utilities"
	sggosdk "github.com/StackGuardian/sg-sdk-go"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	Org          string
	WfgGrp       string
	Preview      bool
	DryRun       bool
	Run          bool
	OutputJson   bool
	PatchPayload string
	Payload      string
}

func NewCreateCmd(c *client.Client) *cobra.Command {
	opts := &RunOptions{}
	// createCmd represents the create command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create new stack",
		Long:  `Create new stack in the specified organization and workflow group.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			opts.Org = cmd.Parent().PersistentFlags().Lookup("org").Value.String()
			opts.WfgGrp = cmd.Parent().PersistentFlags().Lookup("workflow-group").Value.String()
			opts.Payload = args[0]

			payload, err := os.ReadFile(opts.Payload)
			if err != nil {
				cmd.PrintErrln(err)
			}

			var createStackRequest *sggosdk.Stack
			if opts.PatchPayload != "" {
				err := json.Unmarshal(
					[]byte(
						utilities.PatchJSON(string(payload), opts.PatchPayload),
					),
					&createStackRequest)
				if err != nil {
					cmd.Printf("Error during patching Stack payload: %s\n", err)
					os.Exit(-1)
				}
			} else {
				err := json.Unmarshal(
					payload,
					&createStackRequest)
				if err != nil {
					cmd.Printf("Error while unmarshalling Stack payload: %s\n", err)
					os.Exit(-1)
				}
			}
			//Run on create
			if opts.Run {
				createStackRequest.RunOnCreate = sggosdk.Bool(true)
			} else {
				createStackRequest.RunOnCreate = sggosdk.Bool(false)
			}

			// Perform actions based on the set flags
			if err := performPreExecutionFlagChecks(cmd, createStackRequest, opts); err != nil {
				cmd.PrintErrln(err)
				os.Exit(-1)
			}

			response, err := c.Stacks.Create(
				context.Background(),
				opts.Org,
				opts.WfgGrp,
				createStackRequest,
			)
			if err != nil {
				if strings.Contains(err.Error(), "cannot unmarshal") {
					cmd.PrintErrln("Stack was created successfully but an error occured while reading the response.")
					os.Exit(-1)
				}
				cmd.PrintErrln("== Failed To Create Stack ==")
				cmd.PrintErrln(err)
				os.Exit(-1)
			}
			if opts.OutputJson {
				cmd.Println(response)
			}
			cmd.Println("Stack created successfully.")

		},
	}

	// Define the flags for the command

	createCmd.Flags().StringVar(&opts.PatchPayload, "patch-payload", "", "Patch original payload.json input. Add or replace values. Requires valid JSON input.")

	createCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output execution response as json to STDIN.")

	createCmd.Flags().BoolVar(&opts.Preview, "preview", false, "Preview payload content before creating. Execution will not pause.")

	createCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Similar to --preview. But execution will stop, nothing will be created.")

	createCmd.Flags().BoolVar(&opts.Run, "run", false, "Executes the Stack.")

	return createCmd
}

// performPreExecutionFlagChecks performs pre-execution flag checks and returns the payload
func performPreExecutionFlagChecks(cmd *cobra.Command, payload *sggosdk.Stack, opts *RunOptions) error {

	if payload.ResourceName == nil || *payload.ResourceName == "" {
		return errors.New(">> [ERROR] Stack ResourceName is required in object payload, skipping")
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
