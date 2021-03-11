package testutil

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/puppetlabs/leg/k8sutil/pkg/manifest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func doInstallKubernetesManifest(ctx context.Context, cl client.Client, pattern string, patchers ...manifest.PatcherFunc) error {
	files, err := getFixtures(pattern)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("applying manifest %s", file)

		reader, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := ParseApplyKubernetesManifest(ctx, cl, reader, patchers...); err != nil {
			return err
		}
	}

	return nil
}

func doInstall(ctx context.Context, cl client.Client, name string, patchers ...manifest.PatcherFunc) error {
	requested := time.Now()

	pattern := fmt.Sprintf("fixtures/%s/*.yaml", name)
	err := doInstallKubernetesManifest(ctx, cl, pattern, patchers...)
	if err != nil {
		return err
	}

	log.Printf("installed %s after %s", name, time.Now().Sub(requested))
	return nil
}

func doInstallAndWait(ctx context.Context, cl client.Client, namespace, name string, patchers ...manifest.PatcherFunc) error {
	if err := doInstall(ctx, cl, name, patchers...); err != nil {
		return err
	}

	requested := time.Now()

	err := WaitForServicesToBeReady(ctx, cl, namespace)
	if err != nil {
		return err
	}

	log.Printf("waited for services to be ready for %s in %s after %s", name, namespace, time.Now().Sub(requested))
	return nil
}
