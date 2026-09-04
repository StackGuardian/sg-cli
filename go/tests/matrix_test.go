package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/StackGuardian/sg-sdk-go/client"
	"github.com/StackGuardian/sg-sdk-go/option"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recorder is an http.RoundTripper that captures the last request and answers
// with a canned JSON body, so a command's wire format can be asserted exactly.
type recorder struct {
	status int
	resp   string

	method string
	path   string
	query  url.Values
	body   string
	auth   string
	calls  int
}

func (r *recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	r.method, r.path, r.query, r.body, r.auth = req.Method, req.URL.Path, req.URL.Query(), string(body), req.Header.Get("Authorization")
	r.calls++
	status := r.status
	if status == 0 {
		status = http.StatusOK
	}
	resp := r.resp
	if resp == "" {
		resp = "{}"
	}
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(resp)),
	}, nil
}

// matrixCase describes one CLI invocation and the request it must produce.
// "@payload" in args is replaced by a temp file holding payloadJSON.
type matrixCase struct {
	name        string
	ctor        func(*client.Client) *cobra.Command
	args        []string
	payloadJSON string
	method      string
	path        string
	query       url.Values // nil asserts no query string
	body        string     // "" asserts no body
	resp        string
	wantOut     string // substring expected on stdout
}

func runMatrix(t *testing.T, cases []matrixCase) {
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rec := &recorder{resp: tc.resp}
			c := client.NewClient(option.WithApiKey("apikey test-token"), option.WithHTTPClient(&http.Client{Transport: rec}))
			cmd := tc.ctor(c)
			args := make([]string, 0, len(tc.args))
			for _, a := range tc.args {
				if a == "@payload" {
					p := filepath.Join(t.TempDir(), "payload.json")
					require.NoError(t, os.WriteFile(p, []byte(tc.payloadJSON), 0o600))
					a = p
				}
				args = append(args, a)
			}
			cmd.SetArgs(args)
			var out, errOut bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&errOut)
			require.NoError(t, cmd.Execute(), "stderr: %s", errOut.String())
			require.Equal(t, 1, rec.calls, "expected exactly one request")
			assert.Equal(t, tc.method, rec.method, "method")
			assert.Equal(t, tc.path, rec.path, "path")
			if tc.query == nil {
				assert.Empty(t, rec.query, "query")
			} else {
				assert.Equal(t, tc.query, rec.query, "query")
			}
			if tc.body == "" {
				assert.Empty(t, rec.body, "body")
			} else {
				assert.JSONEq(t, tc.body, rec.body, "body")
			}
			assert.Equal(t, "apikey test-token", rec.auth)
			if tc.wantOut != "" {
				assert.Contains(t, out.String(), tc.wantOut)
			} else {
				// read commands print the response as JSON
				var parsed interface{}
				assert.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(out.String())), &parsed), "stdout should be JSON: %s", out.String())
			}
		})
	}
}
