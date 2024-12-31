package moduleregistry

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v64/github"
)

type GitHubRegistry struct {
	Client   *github.Client
	Protocol protocol
}

// New GitHub module registry
func NewGitHubRegistry(config GithubConfig) (Registry, error) {
	var r GitHubRegistry

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.InsecureSkipVerify,
			},
		},
	}

	r.Client = github.NewClient(httpClient)
	r.Protocol = config.Protocol

	if config.TokenSecretName != "" {
		token, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return nil, errors.New("ENV var 'GITHUB_TOKEN' not set")
		}

		r.Client.WithAuthToken(token)
	}

	return r, nil
}

// List all tags for a GitHub registry
// See https://developer.hashicorp.com/terraform/internals/module-registry-protocol#list-available-versions-for-a-specific-module
func (gh GitHubRegistry) Versions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	var allTags []*github.RepositoryTag
	mvs := newModuleVersions()
	group := r.PathValue("group")
	project := r.PathValue("project")
	w.Header().Set("Content-Type", "application/json")

	opt := &github.ListOptions{}
	for {
		tags, resp, err := gh.Client.Repositories.ListTags(ctx, group, project, opt)
		if resp == nil && err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound && err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, v := range allTags {
		mvs.addVersion(v.Name)
	}

	mvs.setSource(fmt.Sprintf("github.com/%s/%s", group, project))

	if err := json.NewEncoder(w).Encode(mvs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Download source code for a specific module version
// See https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version
//
// The module protocl doesn't directly pass the version field as a the ref
// It doesn't want a "v" specified in the HCL but seems to expect tag refs are prefixed with "v"
func (gh GitHubRegistry) Download(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	project := r.PathValue("project")
	version := r.PathValue("version")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-Terraform-Get", fmt.Sprintf("git::%s://github.com/%s/%s?ref=v%s", gh.Protocol, group, project, version))
	w.WriteHeader(http.StatusNoContent)
}
