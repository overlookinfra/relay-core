package spec_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/relay-core/pkg/util/testutil"
	"github.com/puppetlabs/relay-core/pkg/workflow/spec"
	"github.com/stretchr/testify/require"
)

func withStepMetadataServer(t *testing.T, fn func(ts *httptest.Server)) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/step-metadata.json")
		require.NoError(t, err)

		fi, err := f.Stat()
		require.NoError(t, err)

		w.Header().Set("content-type", "application/json")

		http.ServeContent(w, r, "", fi.ModTime(), f)
	}))

	defer ts.Close()

	fn(ts)
}

func TestStepMetadataSchemaRegistry(t *testing.T) {
	var reg spec.SchemaRegistry

	var cases = []struct {
		repo     string
		specFile string
		isValid  bool
	}{
		{
			repo:     "relaysh/kubernetes-step-kubectl",
			specFile: "testdata/kubectl-spec-valid.json",
			isValid:  true,
		},
		{
			repo:     "relaysh/kubernetes-step-kubectl",
			specFile: "testdata/kubectl-spec-invalid.json",
			isValid:  false,
		},
	}

	testutil.WithStepMetadataServer(t, filepath.Join("testdata/step-metadata.json"), func(ts *httptest.Server) {
		u, err := url.Parse(ts.URL)
		require.NoError(t, err)

		stepMetadataReg, err := spec.NewStepMetadataSchemaRegistry(u)
		require.NoError(t, err)

		reg = stepMetadataReg

		for _, c := range cases {
			schema, err := reg.GetByStepRepository(c.repo)
			require.NoError(t, err)

			content, err := ioutil.ReadFile(c.specFile)
			require.NoError(t, err)

			err = schema.Validate(content)

			if c.isValid {
				require.NoError(t, err, errors.Unwrap(err))
			} else {
				require.Error(t, err)
				require.IsType(t, &spec.SchemaValidationError{}, err)
			}
		}
	})
}

func TestStepMetadataSchemaRegistryFetchCaching(t *testing.T) {
	testutil.WithStepMetadataServer(t, filepath.Join("testdata/step-metadata.json"), func(ts *httptest.Server) {
		u, err := url.Parse(ts.URL)
		require.NoError(t, err)

		stepMetadataReg, err := spec.NewStepMetadataSchemaRegistry(u)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, stepMetadataReg.LastResponse.StatusCode)

		_, err = stepMetadataReg.GetByStepRepository("relaysh/kubernetes-step-kubectl")
		require.NoError(t, err)

		require.Equal(t, http.StatusNotModified, stepMetadataReg.LastResponse.StatusCode)
	})
}
