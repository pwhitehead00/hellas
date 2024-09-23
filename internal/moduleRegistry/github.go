package moduleregistry

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v64/github"
	"github.com/pwhitehead00/hellas/internal/models"
)

type GitHubRegistry struct {
	Client   *github.Client
	Protocol string
}

// New GitHub module registry
func NewGitHubRegistry(config GithubConfig) (Registry, error) {
	var r GitHubRegistry

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.InsecureSkipVerify,
			},
		},
	}

	r.Client = github.NewClient(httpClient)
	r.Protocol = config.Protocol

	if config.TokenSecretName != "" {
		token, ok := os.LookupEnv("GitHubToken")
		if !ok {
			return nil, errors.New("ENV var 'GitHubToken' not set")
		}

		r.Client.WithAuthToken(token)
	}

	return r, nil
}

// List all tags for a GitHub registry
func (gh GitHubRegistry) ListVersions(group, project string) (*models.ModuleVersions, error) {
	var allTags []*github.RepositoryTag
	mvs := models.NewModuleVersions()

	opt := &github.ListOptions{}

	for {
		tags, resp, err := gh.Client.Repositories.ListTags(context.Background(), group, project, opt)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, v := range allTags {
		mvs.AddVersion(v.Name)
	}

	return &mvs, nil
}

// Download source code for a specific module version
// See https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version
func (gh *GitHubRegistry) Download(name, provider, version string) string {
	return fmt.Sprintf("git::%s://github.com/%s/%s?ref=v%s", gh.Protocol, name, provider, version)
}
