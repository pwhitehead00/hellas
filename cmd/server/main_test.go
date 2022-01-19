package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEndpoint(method, endpoint string) *httptest.ResponseRecorder {
	os.Setenv("CONFIG", "../../test/github-default.json")
	m := strings.ToUpper(method)
	router := setupRouter("github")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(m, endpoint, nil)
	router.ServeHTTP(w, req)
	return w
}

func TestRootEndpoints(t *testing.T) {
	t.Run("Test Healthcheck Endpoint", func(t *testing.T) {
		w := testEndpoint("GET", "/healthcheck")

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "\"ok\"", w.Body.String())
	})

	t.Run("Test Service Discovery Endpoint", func(t *testing.T) {
		expected := "{\"modules.v1\":\"/v1/modules/\"}"
		w := testEndpoint("GET", "/.well-known/terraform.json")

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expected, w.Body.String())
	})
}

func TestV1Endpoints(t *testing.T) {
	t.Run("Test Download Endpoint", func(t *testing.T) {
		expected := "git::https://github.com/acme/terraform-hapycloud-module?ref=v1.0.0"
		w := testEndpoint("GET", "/v1/modules/acme/module/hapycloud/1.0.0/download")

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, expected, w.Header().Get("X-Terraform-Get"))
	})

	// TODO: Test Versions Endpoint
}
