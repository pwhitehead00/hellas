package gh

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/google/go-github/v40/github"
	"github.com/ironhalo/hellas/models"
)

func MegaGitHub(namespace, provider, name string) models.Modules {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}

	var m models.Modules
	var v models.Versions

	var version []*map[string]string

	ctx := context.Background()
	var allTags []*github.RepositoryTag

	client := github.NewClient(c)
	opt := &github.ListOptions{
		PerPage: 100,
	}

	repo := fmt.Sprintf("terraform-%s-%s", provider, name)
	for {
		tags, resp, err := client.Repositories.ListTags(ctx, namespace, repo, opt)
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
		o := map[string]string{"version": *t.Name}
		version = append(version, &o)
	}
	v.Versions = version
	m.Modules = append(m.Modules, v)

	return m
}
