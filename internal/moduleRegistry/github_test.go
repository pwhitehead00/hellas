package moduleregistry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
				Registries: Registries{
					Github: GithubConfig{
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
