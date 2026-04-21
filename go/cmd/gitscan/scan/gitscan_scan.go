package scan

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// CLI options
// ---------------------------------------------------------------------------

type RunOptions struct {
	Provider        string
	Token           string
	Org             string
	User            string
	MaxRepos        int
	IncludeArchived bool
	IncludeForks    bool
	WfGrp           string
	VCSAuth         string
	ManagedState    bool
	Output          string
	Verbose         bool
	Quiet           bool
}

// ---------------------------------------------------------------------------
// VCS types (normalised repo representation)
// ---------------------------------------------------------------------------

type repo struct {
	ID            string
	Name          string
	FullName      string
	URL           string
	Owner         string
	DefaultBranch string
	IsPrivate     bool
	IsArchived    bool
	IsFork        bool
	Description   string
	Topics        []string
	Provider      string // "GITHUB_COM" | "GITLAB_COM"
}

// ---------------------------------------------------------------------------
// Terraform project (one per directory that contains .tf files)
// ---------------------------------------------------------------------------

type tfProject struct {
	Path        string   // "" for repo root, "infra/vpc" for monorepo subdir
	TFFiles     []string // names of .tf files found
	TFVarsFiles []string // names of .tfvars / .tfvars.json files found
	HasLockfile bool
}

// ---------------------------------------------------------------------------
// sg-payload workflow entry
// ---------------------------------------------------------------------------

type vcsConfig struct {
	IacVCSConfig iacVCSConfig `json:"iacVCSConfig"`
	IacInputData iacInputData `json:"iacInputData"`
}

type iacVCSConfig struct {
	UseMarketplaceTemplate bool         `json:"useMarketplaceTemplate"`
	CustomSource           customSource `json:"customSource"`
}

type customSource struct {
	SourceConfigDestKind string       `json:"sourceConfigDestKind"`
	Config               sourceConfig `json:"config"`
}

type sourceConfig struct {
	Repo             string `json:"repo"`
	Ref              string `json:"ref"`
	IsPrivate        bool   `json:"isPrivate"`
	Auth             string `json:"auth"`
	WorkingDir       string `json:"workingDir"`
	IncludeSubModule bool   `json:"includeSubModule"`
}

type iacInputData struct {
	SchemaType string                 `json:"schemaType"`
	Data       map[string]interface{} `json:"data"`
}

type deploymentPlatformConfig struct {
	Kind   string                 `json:"kind"`
	Config map[string]interface{} `json:"config"`
}

type terraformConfig struct {
	ManagedTerraformState bool   `json:"managedTerraformState"`
	TerraformVersion      string `json:"terraformVersion"`
	ApprovalPreApply      bool   `json:"approvalPreApply"`
	ExtraCLIArgs          string `json:"extraCLIArgs,omitempty"`
}

type runnerConstraints struct {
	Type string `json:"type"`
}

type miniSteps struct {
	WfChaining    wfChaining    `json:"wfChaining"`
	Notifications notifications `json:"notifications"`
}

type wfChaining struct {
	Errored   []interface{} `json:"ERRORED"`
	Completed []interface{} `json:"COMPLETED"`
}

type notifications struct {
	Email emailNotifications `json:"email"`
}

type emailNotifications struct {
	Errored          []interface{} `json:"ERRORED"`
	Completed        []interface{} `json:"COMPLETED"`
	ApprovalRequired []interface{} `json:"APPROVAL_REQUIRED"`
	Cancelled        []interface{} `json:"CANCELLED"`
}

type cliConfiguration struct {
	WorkflowGroup workflowGroup `json:"WorkflowGroup"`
}

type workflowGroup struct {
	Name string `json:"name"`
}

type workflowEntry struct {
	ResourceName             string                     `json:"ResourceName"`
	Description              string                     `json:"Description"`
	Tags                     []string                   `json:"Tags"`
	EnvironmentVariables     []interface{}              `json:"EnvironmentVariables"`
	DeploymentPlatformConfig []deploymentPlatformConfig `json:"DeploymentPlatformConfig"`
	WfType                   string                     `json:"WfType"`
	TerraformConfig          terraformConfig            `json:"TerraformConfig"`
	VCSConfig                vcsConfig                  `json:"VCSConfig"`
	RunnerConstraints        runnerConstraints          `json:"RunnerConstraints"`
	Approvers                []interface{}              `json:"Approvers"`
	MiniSteps                miniSteps                  `json:"MiniSteps"`
	UserSchedules            []interface{}              `json:"UserSchedules"`
	CLIConfiguration         cliConfiguration           `json:"CLIConfiguration"`
}

// ---------------------------------------------------------------------------
// Excluded directories (not Terraform)
// ---------------------------------------------------------------------------

var excludeDirs = map[string]bool{
	".git": true, ".terraform": true, ".terragrunt-cache": true,
	"node_modules": true, "vendor": true, "__pycache__": true,
	".venv": true, "venv": true,
}

// ---------------------------------------------------------------------------
// NewScanCmd
// ---------------------------------------------------------------------------

func NewScanCmd() *cobra.Command {
	opts := &RunOptions{}

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan a GitHub or GitLab organization for Terraform repositories",
		Long: `Scan a GitHub or GitLab organization for Terraform repositories and generate
an sg-payload.json file for bulk workflow creation.

Examples:
  sg-cli git-scan scan --provider github --token ghp_xxx --org my-org
  sg-cli git-scan scan --provider gitlab --token glpat-xxx --org my-group
  sg-cli git-scan scan --provider github --token ghp_xxx --org my-org --max-repos 50 --output export/sg-payload.json`,
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, opts)
		},
	}

	// Required
	cmd.Flags().StringVarP(&opts.Provider, "provider", "p", "", "VCS provider: github or gitlab (required)")
	cmd.Flags().StringVarP(&opts.Token, "token", "t", "", "VCS access token — GitHub PAT or GitLab PAT (required)")
	cmd.MarkFlagRequired("provider")
	cmd.MarkFlagRequired("token")

	// Target
	cmd.Flags().StringVarP(&opts.Org, "org", "o", "", "GitHub organization or GitLab group to scan")
	cmd.Flags().StringVarP(&opts.User, "user", "u", "", "Scan repos for a specific user instead of an org/group")

	// Filtering
	cmd.Flags().IntVarP(&opts.MaxRepos, "max-repos", "m", 0, "Maximum number of repositories to scan (0 = no limit)")
	cmd.Flags().BoolVar(&opts.IncludeArchived, "include-archived", false, "Include archived repositories")
	cmd.Flags().BoolVar(&opts.IncludeForks, "include-forks", false, "Include forked repositories")

	// StackGuardian defaults
	cmd.Flags().StringVar(&opts.WfGrp, "wfgrp", "imported-workflows", "Workflow group name written into CLIConfiguration")
	cmd.Flags().StringVar(&opts.VCSAuth, "vcs-auth", "", "SG VCS integration path (e.g. /integrations/github_com)")
	cmd.Flags().BoolVar(&opts.ManagedState, "managed-state", false, "Enable SG-managed Terraform state")

	// Output
	cmd.Flags().StringVarP(&opts.Output, "output", "O", "sg-payload.json", "Output file path")
	cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Verbose/debug output")
	cmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "Minimal output (warnings and errors only)")

	return cmd
}

// ---------------------------------------------------------------------------
// run — orchestrates discover → scan → transform → write
// ---------------------------------------------------------------------------

func run(cmd *cobra.Command, opts *RunOptions) {
	logger := newLogger(opts)

	// 1. Discover repos
	repos, err := discoverRepos(opts, logger)
	if err != nil {
		cmd.PrintErrln("Error discovering repositories:", err)
		os.Exit(1)
	}
	if len(repos) == 0 {
		cmd.PrintErrln("No repositories found. Check your token, org, and permissions.")
		os.Exit(0)
	}

	// Filter archived / forks
	repos = filterRepos(repos, opts)
	logger.infof("Discovered %d repositories from %s", len(repos), opts.Provider)

	if len(repos) == 0 {
		cmd.PrintErrln("No repositories remain after filtering.")
		os.Exit(0)
	}

	// 2. Scan each repo for Terraform projects
	type repoProjects struct {
		r        repo
		projects []tfProject
	}
	var results []repoProjects

	total := len(repos)
	for idx, r := range repos {
		logger.infof("[%d/%d] Scanning %s...", idx+1, total, r.FullName)
		tree, err := fetchFileTree(r, opts.Token, logger)
		if err != nil {
			logger.warnf("  Could not fetch file tree for %s: %v", r.FullName, err)
			continue
		}
		projects := detectTerraformDirs(tree)
		if len(projects) > 0 {
			logger.infof("  Found %d Terraform project(s) in %s", len(projects), r.FullName)
			results = append(results, repoProjects{r, projects})
		} else {
			logger.debugf("  No Terraform detected in %s", r.FullName)
		}
	}

	if len(results) == 0 {
		cmd.PrintErrln("No Terraform projects found in any repository.")
		os.Exit(0)
	}

	totalProjects := 0
	for _, rp := range results {
		totalProjects += len(rp.projects)
	}
	logger.infof("Found %d Terraform project(s) across %d repo(s)", totalProjects, len(results))

	// 3. Transform to SG payload
	var payload []workflowEntry
	for _, rp := range results {
		for _, proj := range rp.projects {
			wf := buildWorkflow(rp.r, proj, opts)
			payload = append(payload, wf)
		}
	}

	// 4. Write output
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		cmd.PrintErrln("Failed to marshal payload:", err)
		os.Exit(1)
	}

	outputPath := opts.Output
	if dir := filepath.Dir(outputPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			cmd.PrintErrln("Failed to create output directory:", err)
			os.Exit(1)
		}
	}

	if err := os.WriteFile(outputPath, out, 0o644); err != nil {
		cmd.PrintErrln("Failed to write output file:", err)
		os.Exit(1)
	}

	logger.infof("Generated %d workflow(s) → %s", len(payload), outputPath)
	logger.infof("Next step: sg-cli workflow create --bulk --org \"<ORG>\" -- %s", outputPath)
}

// ---------------------------------------------------------------------------
// Repo discovery
// ---------------------------------------------------------------------------

func discoverRepos(opts *RunOptions, logger *logger) ([]repo, error) {
	switch strings.ToLower(opts.Provider) {
	case "github":
		return githubListRepos(opts, logger)
	case "gitlab":
		return gitlabListRepos(opts, logger)
	default:
		return nil, fmt.Errorf("unsupported provider %q — use 'github' or 'gitlab'", opts.Provider)
	}
}

func filterRepos(repos []repo, opts *RunOptions) []repo {
	var out []repo
	for _, r := range repos {
		if !opts.IncludeArchived && r.IsArchived {
			continue
		}
		if !opts.IncludeForks && r.IsFork {
			continue
		}
		out = append(out, r)
	}
	return out
}

// ---------------------------------------------------------------------------
// GitHub client
// ---------------------------------------------------------------------------

const githubAPIURL = "https://api.github.com"

func githubGet(path string, token string, params map[string]string) ([]byte, http.Header, error) {
	u, _ := url.Parse(githubAPIURL + path)
	if params != nil {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("GitHub API %s: HTTP %d: %s", path, resp.StatusCode, string(body))
	}
	return body, resp.Header, nil
}

func githubNextPage(linkHeader string) int {
	// Link: <https://api.github.com/...?page=2>; rel="next", ...
	for _, part := range strings.Split(linkHeader, ",") {
		if strings.Contains(part, `rel="next"`) {
			parts := strings.SplitN(strings.TrimSpace(part), ";", 2)
			rawURL := strings.Trim(parts[0], " <>")
			u, err := url.Parse(rawURL)
			if err != nil {
				continue
			}
			p, err := strconv.Atoi(u.Query().Get("page"))
			if err == nil {
				return p
			}
		}
	}
	return 0
}

func githubListRepos(opts *RunOptions, logger *logger) ([]repo, error) {
	var repos []repo
	page := 1

	for {
		params := map[string]string{
			"per_page": "100",
			"page":     strconv.Itoa(page),
			"type":     "all",
		}

		var path string
		if opts.Org != "" {
			path = "/orgs/" + opts.Org + "/repos"
		} else if opts.User != "" {
			path = "/users/" + opts.User + "/repos"
		} else {
			path = "/user/repos"
			delete(params, "type")
			params["visibility"] = "all"
			params["affiliation"] = "owner,collaborator,organization_member"
		}

		body, headers, err := githubGet(path, opts.Token, params)
		if err != nil {
			return nil, err
		}

		var raw []map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, fmt.Errorf("parsing GitHub response: %w", err)
		}
		if len(raw) == 0 {
			break
		}

		for _, r := range raw {
			repos = append(repos, githubFormatRepo(r))
		}

		if opts.MaxRepos > 0 && len(repos) >= opts.MaxRepos {
			repos = repos[:opts.MaxRepos]
			break
		}

		next := githubNextPage(headers.Get("Link"))
		if next == 0 {
			break
		}
		page = next
	}

	return repos, nil
}

func githubFormatRepo(r map[string]interface{}) repo {
	owner := ""
	if o, ok := r["owner"].(map[string]interface{}); ok {
		owner = stringField(o, "login")
	}
	return repo{
		ID:            fmt.Sprintf("%v", r["id"]),
		Name:          stringField(r, "name"),
		FullName:      stringField(r, "full_name"),
		URL:           stringField(r, "html_url"),
		Owner:         owner,
		DefaultBranch: stringFieldDefault(r, "default_branch", "main"),
		IsPrivate:     boolField(r, "private"),
		IsArchived:    boolField(r, "archived"),
		IsFork:        boolField(r, "fork"),
		Description:   stringField(r, "description"),
		Topics:        stringSliceField(r, "topics"),
		Provider:      "GITHUB_COM",
	}
}

func githubGetFileTree(r repo, token string) ([]string, error) {
	ref := r.DefaultBranch
	if ref == "" {
		ref = "HEAD"
	}
	path := fmt.Sprintf("/repos/%s/%s/git/trees/%s", r.Owner, r.Name, ref)
	body, _, err := githubGet(path+"?recursive=1", token, nil)
	if err != nil {
		// 404 / 409 (empty repo) — not an error, just no files
		if strings.Contains(err.Error(), "HTTP 404") || strings.Contains(err.Error(), "HTTP 409") {
			return nil, nil
		}
		return nil, err
	}

	var result struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var files []string
	for _, item := range result.Tree {
		if item.Type == "blob" {
			files = append(files, item.Path)
		}
	}
	return files, nil
}

// ---------------------------------------------------------------------------
// GitLab client
// ---------------------------------------------------------------------------

const gitlabAPIURL = "https://gitlab.com/api/v4"

func gitlabGet(path string, token string, params map[string]string) ([]byte, http.Header, error) {
	baseURL := os.Getenv("GITLAB_API_URL")
	if baseURL == "" {
		baseURL = gitlabAPIURL
	}

	u, _ := url.Parse(baseURL + path)
	if params != nil {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("GitLab API %s: HTTP %d: %s", path, resp.StatusCode, string(body))
	}
	return body, resp.Header, nil
}

func gitlabNextPage(headers http.Header) int {
	np := strings.TrimSpace(headers.Get("X-Next-Page"))
	if np == "" {
		return 0
	}
	p, err := strconv.Atoi(np)
	if err != nil {
		return 0
	}
	return p
}

func gitlabListRepos(opts *RunOptions, logger *logger) ([]repo, error) {
	var repos []repo
	page := 1

	for {
		params := map[string]string{
			"per_page":      "100",
			"page":          strconv.Itoa(page),
			"order_by":      "last_activity_at",
			"sort":          "desc",
		}

		var path string
		if opts.Org != "" {
			encoded := url.PathEscape(opts.Org)
			params["include_subgroups"] = "true"
			path = "/groups/" + encoded + "/projects"
		} else if opts.User != "" {
			path = "/users/" + opts.User + "/projects"
		} else {
			params["membership"] = "true"
			path = "/projects"
		}

		body, headers, err := gitlabGet(path, opts.Token, params)
		if err != nil {
			return nil, err
		}

		var raw []map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, fmt.Errorf("parsing GitLab response: %w", err)
		}
		if len(raw) == 0 {
			break
		}

		for _, p := range raw {
			repos = append(repos, gitlabFormatRepo(p))
		}

		if opts.MaxRepos > 0 && len(repos) >= opts.MaxRepos {
			repos = repos[:opts.MaxRepos]
			break
		}

		next := gitlabNextPage(headers)
		if next == 0 {
			break
		}
		page = next
	}

	return repos, nil
}

func gitlabFormatRepo(p map[string]interface{}) repo {
	pathWithNS := stringField(p, "path_with_namespace")
	parts := strings.SplitN(pathWithNS, "/", 2)
	owner := ""
	if len(parts) > 1 {
		owner = parts[0]
	}

	// topics may be under "topics" or legacy "tag_list"
	topics := stringSliceField(p, "topics")
	if len(topics) == 0 {
		topics = stringSliceField(p, "tag_list")
	}

	isFork := false
	if fp, ok := p["forked_from_project"]; ok && fp != nil {
		isFork = true
	}

	return repo{
		ID:            fmt.Sprintf("%v", p["id"]),
		Name:          stringField(p, "name"),
		FullName:      pathWithNS,
		URL:           stringField(p, "web_url"),
		Owner:         owner,
		DefaultBranch: stringFieldDefault(p, "default_branch", "main"),
		IsPrivate:     stringField(p, "visibility") == "private",
		IsArchived:    boolField(p, "archived"),
		IsFork:        isFork,
		Description:   stringField(p, "description"),
		Topics:        topics,
		Provider:      "GITLAB_COM",
	}
}

func gitlabGetFileTree(r repo, token string) ([]string, error) {
	ref := r.DefaultBranch
	if ref == "" {
		ref = "HEAD"
	}

	var files []string
	page := 1

	for {
		params := map[string]string{
			"ref":       ref,
			"recursive": "true",
			"per_page":  "100",
			"page":      strconv.Itoa(page),
		}
		body, headers, err := gitlabGet("/projects/"+r.ID+"/repository/tree", token, params)
		if err != nil {
			if strings.Contains(err.Error(), "HTTP 404") || strings.Contains(err.Error(), "HTTP 409") {
				return nil, nil
			}
			return nil, err
		}

		var items []map[string]interface{}
		if err := json.Unmarshal(body, &items); err != nil {
			return nil, err
		}
		if len(items) == 0 {
			break
		}

		for _, item := range items {
			if stringField(item, "type") == "blob" {
				files = append(files, stringField(item, "path"))
			}
		}

		next := gitlabNextPage(headers)
		if next == 0 {
			break
		}
		page = next
	}

	return files, nil
}

// ---------------------------------------------------------------------------
// File tree dispatch
// ---------------------------------------------------------------------------

func fetchFileTree(r repo, token string, logger *logger) ([]string, error) {
	switch r.Provider {
	case "GITHUB_COM":
		return githubGetFileTree(r, token)
	case "GITLAB_COM":
		return gitlabGetFileTree(r, token)
	default:
		return nil, fmt.Errorf("unknown provider %q", r.Provider)
	}
}

// ---------------------------------------------------------------------------
// Terraform detection (file-tree based, no cloning)
// ---------------------------------------------------------------------------

func detectTerraformDirs(fileTree []string) []tfProject {
	type dirEntry struct {
		tfFiles     []string
		tfvarsFiles []string
		hasLockfile bool
	}

	dirs := map[string]*dirEntry{}

	for _, filePath := range fileTree {
		parts := strings.Split(filePath, "/")

		// skip excluded directories
		excluded := false
		for _, p := range parts {
			if excludeDirs[p] {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		name := parts[len(parts)-1]
		dirPath := ""
		if len(parts) > 1 {
			dirPath = strings.Join(parts[:len(parts)-1], "/")
		}

		if _, ok := dirs[dirPath]; !ok {
			dirs[dirPath] = &dirEntry{}
		}
		e := dirs[dirPath]

		switch {
		case strings.HasSuffix(name, ".tf") || strings.HasSuffix(name, ".tf.json"):
			e.tfFiles = append(e.tfFiles, name)
		case strings.HasSuffix(name, ".tfvars") || strings.HasSuffix(name, ".tfvars.json"):
			e.tfvarsFiles = append(e.tfvarsFiles, name)
		case name == ".terraform.lock.hcl":
			e.hasLockfile = true
		}
	}

	var projects []tfProject
	for dirPath, e := range dirs {
		if len(e.tfFiles) == 0 {
			continue
		}
		projects = append(projects, tfProject{
			Path:        dirPath,
			TFFiles:     e.tfFiles,
			TFVarsFiles: e.tfvarsFiles,
			HasLockfile: e.hasLockfile,
		})
	}
	return projects
}

// ---------------------------------------------------------------------------
// Transform: repo + project → workflowEntry
// ---------------------------------------------------------------------------

func buildWorkflow(r repo, proj tfProject, opts *RunOptions) workflowEntry {
	// Resource name
	resourceName := r.Name
	if proj.Path != "" && proj.Path != "." {
		resourceName = r.Name + "-" + strings.ReplaceAll(proj.Path, "/", "-")
	}

	// Description
	description := r.Description
	if description == "" {
		description = "Workflow for " + r.FullName
	}

	// Tags
	tags := append([]string{}, r.Topics...)
	tags = append(tags, "terraform")
	tags = dedupe(tags)

	// VCS auth
	auth := opts.VCSAuth
	if !r.IsPrivate {
		auth = ""
	}

	// Working dir
	workingDir := proj.Path
	if workingDir == "." {
		workingDir = ""
	}

	// Deployment platform config — placeholder; user fills in integration ID.
	// API requires at least one entry — emit a placeholder so the field is never empty.
	deployConfig := []deploymentPlatformConfig{
		{
			Kind: "TERRAFORM_OTHER",
			Config: map[string]interface{}{
				"integrationId": "PLEASE_CONFIGURE",
			},
		},
	}

	// TFVars extra CLI args
	extraCLIArgs := ""
	if len(proj.TFVarsFiles) > 0 {
		extraCLIArgs = "-var-file=" + proj.TFVarsFiles[0]
	}

	tfConfig := terraformConfig{
		ManagedTerraformState: opts.ManagedState,
		TerraformVersion:      "1.5.0", // default; no HCL parsing without clone
		ApprovalPreApply:      true,
		ExtraCLIArgs:          extraCLIArgs,
	}

	wf := workflowEntry{
		ResourceName:             resourceName,
		Description:              description,
		Tags:                     tags,
		EnvironmentVariables:     []interface{}{},
		DeploymentPlatformConfig: deployConfig,
		WfType:                   "TERRAFORM",
		TerraformConfig:          tfConfig,
		VCSConfig: vcsConfig{
			IacVCSConfig: iacVCSConfig{
				UseMarketplaceTemplate: false,
				CustomSource: customSource{
					SourceConfigDestKind: r.Provider,
					Config: sourceConfig{
						Repo:             r.URL,
						Ref:              r.DefaultBranch,
						IsPrivate:        r.IsPrivate,
						Auth:             auth,
						WorkingDir:       workingDir,
						IncludeSubModule: false,
					},
				},
			},
			IacInputData: iacInputData{
				SchemaType: "RAW_JSON",
				Data:       map[string]interface{}{},
			},
		},
		RunnerConstraints: runnerConstraints{Type: "shared"},
		Approvers:         []interface{}{},
		MiniSteps: miniSteps{
			WfChaining: wfChaining{
				Errored:   []interface{}{},
				Completed: []interface{}{},
			},
			Notifications: notifications{
				Email: emailNotifications{
					Errored:          []interface{}{},
					Completed:        []interface{}{},
					ApprovalRequired: []interface{}{},
					Cancelled:        []interface{}{},
				},
			},
		},
		UserSchedules: []interface{}{},
		CLIConfiguration: cliConfiguration{
			WorkflowGroup: workflowGroup{Name: opts.WfGrp},
		},
	}

	return wf
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func dedupe(s []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func stringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func stringFieldDefault(m map[string]interface{}, key, def string) string {
	s := stringField(m, key)
	if s == "" {
		return def
	}
	return s
}

func boolField(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func stringSliceField(m map[string]interface{}, key string) []string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var out []string
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Logger
// ---------------------------------------------------------------------------

type logger struct {
	verbose bool
	quiet   bool
	log     *log.Logger
}

func newLogger(opts *RunOptions) *logger {
	return &logger{
		verbose: opts.Verbose,
		quiet:   opts.Quiet,
		log:     log.New(os.Stderr, "", 0),
	}
}

func (l *logger) infof(format string, args ...interface{}) {
	if !l.quiet {
		l.log.Printf("[INFO]  "+format, args...)
	}
}

func (l *logger) warnf(format string, args ...interface{}) {
	l.log.Printf("[WARN]  "+format, args...)
}

func (l *logger) debugf(format string, args ...interface{}) {
	if l.verbose {
		l.log.Printf("[DEBUG] "+format, args...)
	}
}
