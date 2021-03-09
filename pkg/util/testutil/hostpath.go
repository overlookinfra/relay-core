// https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/tree/master/examples/hostpath-provisioner

package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InstallHostpathProvisioner(t *testing.T, ctx context.Context, cl client.Client) {
	require.NoError(t, doInstall(ctx, cl, "hostpath"))
}
