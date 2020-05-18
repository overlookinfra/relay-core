package testutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/nebula-tasks/pkg/util/testutil"
)

func TestInstallKnativeServing(t *testing.T) {
	testutil.WithEndToEndEnvironment(t, func(e2e *testutil.EndToEndEnvironment) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		testutil.InstallKnativeServing(t, ctx, e2e.ControllerRuntimeClient)
	})
}
