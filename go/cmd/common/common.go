// Package common holds the helpers shared by every sg-cli command: client
// construction, JSON output, flag access and payload handling.
package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/StackGuardian/sg-cli/utilities"
	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/option"
	"github.com/spf13/cobra"
)

const (
	DefaultBaseURL = "https://api.app.stackguardian.io"
	DashboardURL   = "https://app.stackguardian.io/orchestrator"
)

// capture records the most recent response body so commands can print the
// API's JSON verbatim: re-encoding the SDK's typed structs would drop any
// field the SDK does not declare.
type capture struct {
	next http.RoundTripper
	mu   sync.Mutex
	body []byte
}

func (c *capture) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))
	c.mu.Lock()
	c.body = body
	c.mu.Unlock()
	return resp, nil
}

var lastResponse = &capture{next: http.DefaultTransport}

// NewClient builds the SDK client from SG_API_TOKEN and SG_BASE_URL with the
// response capture installed. Both sgu_ (personal) and sgo_ (org) tokens work.
func NewClient() *client.Client {
	apiKey := strings.TrimSpace(os.Getenv("SG_API_TOKEN"))
	if !strings.HasPrefix(strings.ToLower(apiKey), "apikey ") {
		apiKey = "apikey " + apiKey
	}
	baseURL := os.Getenv("SG_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return client.NewClient(
		option.WithApiKey(apiKey),
		option.WithBaseURL(strings.TrimRight(baseURL, "/")),
		option.WithHTTPClient(&http.Client{Transport: lastResponse}),
	)
}

// PrintJSON writes the API response as indented JSON: the captured raw body
// when the client was built by NewClient, otherwise the encoded value.
func PrintJSON(cmd *cobra.Command, response interface{}) {
	lastResponse.mu.Lock()
	body := lastResponse.body
	lastResponse.mu.Unlock()
	var out bytes.Buffer
	if len(bytes.TrimSpace(body)) > 0 && json.Indent(&out, body, "", "  ") == nil {
		cmd.Println(out.String())
		return
	}
	encoded, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		cmd.Println(response)
		return
	}
	cmd.Println(string(encoded))
}

// Org returns the value of the --org flag, searching parent commands.
func Org(cmd *cobra.Command) string {
	return String(cmd, "org")
}

// String returns a string flag defined on cmd or inherited from a parent.
func String(cmd *cobra.Command, name string) string {
	if f := cmd.Flag(name); f != nil {
		return f.Value.String()
	}
	return ""
}

// Bool returns a bool flag defined on cmd or inherited from a parent.
func Bool(cmd *cobra.Command, name string) bool {
	if f := cmd.Flag(name); f != nil {
		return f.Value.String() == "true"
	}
	return false
}

// Optional flag helpers: they return nil when the flag was not set, so a
// default never overrides the API's own default.

func OptString(cmd *cobra.Command, name string) *string {
	if f := cmd.Flags().Lookup(name); f == nil || !f.Changed {
		return nil
	}
	v, _ := cmd.Flags().GetString(name)
	return &v
}

func OptInt(cmd *cobra.Command, name string) *int {
	if f := cmd.Flags().Lookup(name); f == nil || !f.Changed {
		return nil
	}
	v, _ := cmd.Flags().GetInt(name)
	return &v
}

func OptBool(cmd *cobra.Command, name string) *bool {
	if f := cmd.Flags().Lookup(name); f == nil || !f.Changed {
		return nil
	}
	v, _ := cmd.Flags().GetBool(name)
	return &v
}

func OptStringSlice(cmd *cobra.Command, name string) []string {
	if f := cmd.Flags().Lookup(name); f == nil || !f.Changed {
		return nil
	}
	v, _ := cmd.Flags().GetStringSlice(name)
	return v
}

// AddOrgFlag registers the required --org persistent flag on a command group.
func AddOrgFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String("org", "", "The organization name on Stackguardian platform.")
	_ = cmd.MarkPersistentFlagRequired("org")
}

// AddPayloadFlags registers the flags that accompany a JSON payload file.
func AddPayloadFlags(cmd *cobra.Command) {
	cmd.Flags().String("patch-payload", "", "Patch original payload.json input. Add or replace values. Requires valid JSON input.")
	cmd.Flags().Bool("preview", false, "Preview payload content before sending. Execution will not pause.")
	cmd.Flags().Bool("dry-run", false, "Similar to --preview. But execution will stop, nothing will be sent.")
}

// AddOutputJSONFlag registers --output-json on commands that print a message by default.
func AddOutputJSONFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("output-json", false, "Output execution response as json to STDOUT.")
}

// ErrDryRun is returned by ReadPayload after printing the payload for --dry-run.
var ErrDryRun = errors.New("dry run: nothing was sent")

// ReadPayload reads the JSON payload file given as the positional argument
// (use "-" for stdin), applies --patch-payload, honours --preview/--dry-run
// and decodes the result into dst.
func ReadPayload(cmd *cobra.Command, args []string, dst interface{}) error {
	if len(args) != 1 {
		return errors.New("a JSON payload file is required after --")
	}
	var payload []byte
	var err error
	if args[0] == "-" {
		payload, err = io.ReadAll(cmd.InOrStdin())
	} else {
		payload, err = os.ReadFile(args[0])
	}
	if err != nil {
		return fmt.Errorf("read payload: %w", err)
	}
	if patch := String(cmd, "patch-payload"); patch != "" {
		payload = []byte(utilities.PatchJSON(string(payload), patch))
	}
	if Bool(cmd, "preview") || Bool(cmd, "dry-run") {
		var pretty bytes.Buffer
		if json.Indent(&pretty, payload, "", "    ") == nil {
			cmd.Println(pretty.String())
		} else {
			cmd.Println(string(payload))
		}
		if Bool(cmd, "dry-run") {
			return ErrDryRun
		}
	}
	if err := json.Unmarshal(payload, dst); err != nil {
		return fmt.Errorf("invalid JSON payload: %w", err)
	}
	return nil
}

// Done prints the response when --output-json is set, then the message.
func Done(cmd *cobra.Command, response interface{}, message string) {
	if Bool(cmd, "output-json") {
		PrintJSON(cmd, response)
	}
	cmd.Println(message)
}

// Require returns an error naming the first of the given flags that is empty.
func Require(cmd *cobra.Command, names ...string) error {
	for _, n := range names {
		if String(cmd, n) == "" {
			return fmt.Errorf("required flag \"%s\" not set", n)
		}
	}
	return nil
}

// CheckEnum validates that value is one of allowed.
func CheckEnum(flag, value string, allowed ...string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	return fmt.Errorf("invalid value %q for --%s: must be one of %s", value, flag, strings.Join(allowed, ", "))
}
