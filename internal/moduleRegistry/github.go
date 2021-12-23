package moduleregistry

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/google/go-github/v40/github"
	"github.com/ironhalo/hellas/models"
)

type GitHubClient struct {
	Client *github.Client
}

func NewGitHubClient() Registry {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}
	client := github.NewClient(c)

	return &GitHubClient{
		Client: client,
	}
}

func (gh *GitHubClient) Versions(namespace, name, provider string) models.ModuleVersions {
	var m models.ModuleVersions
	var versions []*models.ModuleVersion
	var allTags []*github.RepositoryTag
	ctx := context.Background()

	opt := &github.ListOptions{
		PerPage: 100,
	}

	repo := fmt.Sprintf("terraform-%s-%s", provider, name)
	for {
		tags, resp, err := gh.Client.Repositories.ListTags(ctx, namespace, repo, opt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, t := range allTags {
		o := models.ModuleVersion{
			Version: *t.Name,
		}
		versions = append(versions, &o)
	}

	mpv := models.ModuleProviderVersions{
		Source:   repo,
		Versions: versions,
	}

	m.Modules = append(m.Modules, &mpv)

	return m
}

func (gh *GitHubClient) Download(namespace, name, provider, version string) (source string) {
	source = fmt.Sprintf("git::https://github.com/%s/terraform-%s-%s?ref=v%s", namespace, provider, name, version)
	return
}
