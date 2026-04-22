package create

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/StackGuardian/sg-cli/cmd/output"
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

	var createCmd = &cobra.Command{
		Use:   "create [payload.json]",
		Short: "Create a new stack",
		Long:  "Create a new stack from a JSON payload file.",
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

			var createStackRequest *sggosdk.Stack
			if opts.PatchPayload != "" {
				if err := json.Unmarshal([]byte(utilities.PatchJSON(string(payload), opts.PatchPayload)), &createStackRequest); err != nil {
					output.Error("Error patching stack payload: " + err.Error())
					os.Exit(1)
				}
			} else {
				if err := json.Unmarshal(payload, &createStackRequest); err != nil {
					output.Error("Error reading stack payload: " + err.Error())
					os.Exit(1)
				}
			}

			createStackRequest.RunOnCreate = sggosdk.Bool(opts.Run)

			if err := performPreExecutionFlagChecks(cmd, createStackRequest, opts); err != nil {
				output.Error(err.Error())
				os.Exit(1)
			}

			var response interface{}
			err = output.WithSpinner("Creating stack...", func() error {
				var e error
				response, e = c.Stacks.CreateStack(context.Background(), opts.Org, opts.WfgGrp, createStackRequest)
				return e
			})
			if err != nil {
				if strings.Contains(err.Error(), "cannot unmarshal") {
					output.Success("Stack created successfully.")
					output.Warning("Could not parse the API response JSON.")
					os.Exit(0)
				}
				output.Error("Failed to create stack: " + err.Error())
				os.Exit(1)
			}

			if opts.OutputJson {
				cmd.Println(response)
			}
			output.Success("Stack created successfully.")
		},
	}

	createCmd.Flags().StringVar(&opts.PatchPayload, "patch-payload", "", "Merge a JSON patch over the payload before applying.")
	createCmd.Flags().BoolVar(&opts.OutputJson, "output-json", false, "Output API response as JSON.")
	createCmd.Flags().BoolVar(&opts.Preview, "preview", false, "Preview the payload before creating (execution continues).")
	createCmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the payload and exit without creating.")
	createCmd.Flags().BoolVar(&opts.Run, "run", false, "Trigger a run immediately after creation.")

	return createCmd
}

// performPreExecutionFlagChecks validates and optionally previews the stack payload.
func performPreExecutionFlagChecks(cmd *cobra.Command, payload *sggosdk.Stack, opts *RunOptions) error {
	if payload.ResourceName == nil || payload.ResourceName.Value == "" {
		return errors.New("stack ResourceName is required in the payload")
	}

	if opts.DryRun || opts.Preview {
		requestJson, err := json.MarshalIndent(payload, "", "    ")
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
		}
		output.Section("Payload Preview")
		cmd.Println(string(requestJson))
		if opts.DryRun {
			os.Exit(0)
		}
	}
	return nil
}
