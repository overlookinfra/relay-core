package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/puppetlabs/relay-core/pkg/metadataapi/opt"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/sample"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/server/api"
	"github.com/puppetlabs/relay-core/pkg/spec"
	"github.com/stretchr/testify/require"
)

func TestGetEnvironment(t *testing.T) {
	ctx := context.Background()

	tokenGenerator, err := sample.NewHS256TokenGenerator(nil)
	require.NoError(t, err)

	sc := &opt.SampleConfig{
		Secrets: map[string]string{
			"test-secret-key": "test-secret-value",
		},
		Runs: map[string]*opt.SampleConfigRun{
			"test": {
				Steps: map[string]*opt.SampleConfigStep{
					"previous-task": {
						Outputs: map[string]interface{}{
							"test-output-key": "test-output-value",
						},
					},
					"current-task": {
						Spec: opt.SampleConfigSpec{
							"superSecret": spec.YAMLTree{
								Tree: "${secrets.test-secret-key}",
							},
							"superOutput": spec.YAMLTree{
								Tree: "${outputs.previous-task.test-output-key}",
							},
						},
						Env: opt.SampleConfigEnvironment{
							"test-environment-variable-from-secret": spec.YAMLTree{
								Tree: "${secrets.test-secret-key}",
							},
							"test-environment-variable-from-output": spec.YAMLTree{
								Tree: "${outputs.previous-task.test-output-key}",
							},
						},
					},
				},
			},
		},
	}

	tokenMap := tokenGenerator.GenerateAll(ctx, sc)

	currentTaskToken, found := tokenMap.ForStep("test", "current-task")
	require.True(t, found)

	h := api.NewHandler(sample.NewAuthenticator(sc, tokenGenerator.Key()))

	req, err := http.NewRequest(http.MethodGet, "/spec", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+currentTaskToken)

	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Result().StatusCode)

	var r api.GetSpecResponseEnvelope
	require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&r))
	require.Equal(t, map[string]interface{}{
		"superSecret": "test-secret-value",
		"superOutput": "test-output-value",
	}, r.Value.Data)
	require.True(t, r.Complete)

	req, err = http.NewRequest(http.MethodGet, "/environment", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+currentTaskToken)

	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Result().StatusCode)

	require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&r))
	require.Equal(t, map[string]interface{}{
		"test-environment-variable-from-secret": "test-secret-value",
		"test-environment-variable-from-output": "test-output-value",
	}, r.Value.Data)
	require.True(t, r.Complete)

	req, err = http.NewRequest(http.MethodGet, "/environment/test-environment-variable-from-secret", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+currentTaskToken)

	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Result().StatusCode)

	require.NoError(t, json.NewDecoder(resp.Result().Body).Decode(&r))
	require.Equal(t, "test-secret-value", r.Value.Data)

	req, err = http.NewRequest(http.MethodGet, "/environment/test-environment-variable-from-output", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+currentTaskToken)

	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Result().StatusCode)

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&r))
	require.Equal(t, "test-output-value", r.Value.Data)
}
