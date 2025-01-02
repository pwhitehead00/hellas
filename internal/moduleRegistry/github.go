package moduleregistry

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v64/github"
	"github.com/pwhitehead00/hellas/internal/logging"
)

type gitHubRegistry struct {
	client   *github.Client
	protocol protocol
}

// New GitHub module registry
func NewGitHubRegistry(httpClient *http.Client, config githubConfig) (registry, error) {
	var r gitHubRegistry
	r.client = github.NewClient(httpClient)
	r.protocol = config.Protocol

	if config.TokenSecretName != "" {
		token, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return nil, errors.New("ENV var 'GITHUB_TOKEN' not set")
		}

		r.client.WithAuthToken(token)
	}

	return r, nil
}

// List all tags for a GitHub registry
// See https://developer.hashicorp.com/terraform/internals/module-registry-protocol#list-available-versions-for-a-specific-module
func (gh gitHubRegistry) Versions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var allTags []*github.RepositoryTag
		mvs := newModuleVersions()
		group := r.PathValue("group")
		project := r.PathValue("project")
		w.Header().Set("Content-Type", "application/json")
		source := fmt.Sprintf("github.com/%s/%s", group, project)

		log := logging.Log.With("handler", "versions", "repository-type", "github", "source", source)

		opt := &github.ListOptions{}
		for {
			tags, resp, err := gh.client.Repositories.ListTags(r.Context(), group, project, opt)
			if resp == nil && err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Error("failed to list tags", "error", err.Error())
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound && err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				log.Error("project not found", "error", err.Error())
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

		mvs.setSource(source)

		if err := json.NewEncoder(w).Encode(mvs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("found tags successfully", "found", len(allTags))
	}
}

// Download source code for a specific module version
// See https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version
//
// The module protocl doesn't directly pass the version field as a the ref
// It doesn't want a "v" specified in the HCL but seems to expect tag refs are prefixed with "v"
func (gh gitHubRegistry) Download() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.PathValue("group")
		project := r.PathValue("project")
		version := r.PathValue("version")
		source := fmt.Sprintf("git::%s://github.com/%s/%s?ref=v%s", gh.protocol, group, project, version)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("X-Terraform-Get", source)
		w.WriteHeader(http.StatusNoContent)

		log := logging.Log.With("handler", "download", "repository-type", "github")
		log.Info("X-Terraform-Get", "value", source)
	}
}
