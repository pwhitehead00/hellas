package moduleregistry

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModuleRegistry(t *testing.T) {
	testCases := []struct {
		desc  string
		input Config
		err   error
	}{
		{
			desc: "Valid GitHub Config",
			input: Config{
				Registries: registries{
					Github: githubConfig{
						Protocol:           "https",
						InsecureSkipVerify: true,
						Enabled:            true,
					},
				},
			},
			err: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := NewModuleRegistry(tC.input)
			assert.Equal(t, tC.err, err)
		})
	}
}

func TestDiscoveryEndpoint(t *testing.T) {
	testCases := []struct {
		desc       string
		statusCode int
		body       string
	}{
		{
			desc:       "successful response",
			statusCode: http.StatusOK,
			body:       "{\"modules.v1\":\"/v1/modules/\"}\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/.well-known/terraform.json", nil)
			w := httptest.NewRecorder()
			discovery(w, req)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				assert.Fail(t, err.Error(), "unexpected error reading response body")
			}
			defer resp.Body.Close()

			assert.Equal(t, tC.body, string(body))
			assert.Equal(t, tC.statusCode, resp.StatusCode)
		})
	}
}

func TestHealthCheck(t *testing.T) {
	testCases := []struct {
		desc       string
		statusCode int
		body       string
	}{
		{
			desc:       "successful response",
			statusCode: http.StatusOK,
			body:       "\"ok\"\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
			w := httptest.NewRecorder()
			healthCheck(w, req)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				assert.Fail(t, err.Error(), "unexpected error reading response body")
			}
			defer resp.Body.Close()

			assert.Equal(t, tC.body, string(body))
			assert.Equal(t, tC.statusCode, resp.StatusCode)
		})
	}
}
