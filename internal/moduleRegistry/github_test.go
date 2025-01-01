package moduleregistry

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-github/v64/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestGitHubDownload(t *testing.T) {
	testCases := []struct {
		desc            string
		config          Config
		uri             string
		contentType     string
		terraformHeader string
		status          int
	}{
		{
			desc: "successful http protocol",
			config: Config{
				Registries: registries{
					Github: githubConfig{
						Enabled:  true,
						Protocol: protocolHTTPS,
					},
				},
			},
			uri:             "/v1/modules/my-group/my-project/github/1.2.3/download",
			contentType:     "application/json",
			terraformHeader: "git::https://github.com/my-group/my-project?ref=v1.2.3",
			status:          http.StatusNoContent,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := NewModuleRegistry(tC.config)
			if err != nil {
				assert.Fail(t, err.Error(), "unexpected error: newModuleRegistry")
			}

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := ts.Client().Get(fmt.Sprintf("%s%s", ts.URL, tC.uri))

			if err != nil {
				assert.Fail(t, err.Error(), "unexpected error: client.Get")
			}

			assert.Equal(t, tC.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tC.terraformHeader, resp.Header.Get("X-Terraform-Get"))
			assert.Equal(t, tC.status, resp.StatusCode)
		})
	}
}

func TestGitHubVersions(t *testing.T) {
	testCases := []struct {
		desc     string
		client   *http.Client
		expected string
		status   int
	}{
		{
			desc: "successful get repository tags",
			client: mock.NewMockedHTTPClient(
				mock.WithRequestMatchPages(
					mock.GetReposTagsByOwnerByRepo, []github.RepositoryTag{
						{
							Name: github.String("v1.2.3"),
						},
						{
							Name: github.String("v1.2.4"),
						},
					},
					mock.GetReposTagsByOwnerByRepo, []github.RepositoryTag{
						{
							Name: github.String("v1.2.5"),
						},
						{
							Name: github.String("v1.3.0"),
						},
					},
				),
			),
			expected: `{"modules":[{"source":"github.com/my-group/my-project","versions":[{"version":"v1.2.3"},{"version":"v1.2.4"},{"version":"v1.2.5"},{"version":"v1.3.0"}]}]}`,
			status:   http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			config := githubConfig{
				Protocol: protocolHTTPS,
				Enabled:  true,
			}
			r, err := NewGitHubRegistry(tC.client, config)
			if err != nil {
				assert.Fail(t, err.Error(), "unexpected error creating NewGitHubRegistry")
			}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /v1/modules/{group}/{project}/github/versions", r.Versions())

			ts := httptest.NewTLSServer(mux)
			defer ts.Close()

			resp, err := ts.Client().Get(fmt.Sprintf("%s%s", ts.URL, "/v1/modules/my-group/my-project/github/versions"))
			if err != nil {
				assert.Fail(t, err.Error(), "unexpected error on get")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				assert.Fail(t, err.Error(), "unexpected error reading response body")
			}

			assert.Equal(t, tC.expected, strings.TrimSpace(string(body)))
			assert.Equal(t, tC.status, resp.StatusCode)
		})
	}
}
