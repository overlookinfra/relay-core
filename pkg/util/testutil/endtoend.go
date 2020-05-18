package testutil

import (
	"context"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/puppetlabs/nebula-tasks/pkg/obj"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type EndToEndEnvironment struct {
	RESTConfig              *rest.Config
	RESTMapper              meta.RESTMapper
	ControllerRuntimeClient client.Client
	Interface               kubernetes.Interface
	TektonInterface         tekton.Interface
}

func (e *EndToEndEnvironment) WithTestNamespace(t *testing.T, ctx context.Context, fn func(ns *corev1.Namespace)) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("relay-e2e-%s", uuid.New()),
			Labels: map[string]string{
				"testing.relay.sh/harness": "end-to-end",
			},
		},
	}
	require.NoError(t, e.ControllerRuntimeClient.Create(ctx, ns))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		assert.NoError(t, e.ControllerRuntimeClient.Delete(ctx, ns))
	}()

	fn(ns)
}

type EndToEndEnvironmentOption func(ctx context.Context, e *EndToEndEnvironment) error

func EndToEndEnvironmentWithTekton(ctx context.Context, e *EndToEndEnvironment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	return doInstallTektonPipeline(ctx, e.ControllerRuntimeClient)
}

func EndToEndEnvironmentWithKnative(ctx context.Context, e *EndToEndEnvironment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	return doInstallKnativeServing(ctx, e.ControllerRuntimeClient)
}

func EndToEndEnvironmentWithAmbassador(ctx context.Context, e *EndToEndEnvironment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	return doInstallAmbassador(ctx, e.ControllerRuntimeClient, e.RESTMapper)
}

var _ EndToEndEnvironmentOption = EndToEndEnvironmentWithTekton
var _ EndToEndEnvironmentOption = EndToEndEnvironmentWithKnative

func doEndToEndEnvironment(fn func(e *EndToEndEnvironment), opts ...EndToEndEnvironmentOption) (bool, error) {
	// We'll allow 30 minutes to attach the environment and run the test. This
	// should give us time to acquire the exclusive lock but puts a cap on any
	// underlying issues causing timeouts in CI.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	viper.SetEnvPrefix("relay_test_e2e")
	viper.AutomaticEnv()

	viper.SetDefault("label_nodes", false)
	viper.SetDefault("disabled", false)
	viper.SetDefault("install_environment", true)

	if viper.GetBool("disabled") {
		return false, nil
	}

	// Don't inherit from $KUBECONFIG so people don't accidentally run these
	// tests against a cluster they care about.
	kubeconfigs := strings.TrimSpace(viper.GetString("kubeconfig"))
	if kubeconfigs == "" {
		return true, fmt.Errorf("end-to-end tests require the RELAY_TEST_E2E_KUBECONFIG environment variable to be set to the path of a valid Kubeconfig")
	}

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{
			Precedence:       filepath.SplitList(kubeconfigs),
			WarnIfAllMissing: true,
		},
		&clientcmd.ConfigOverrides{
			CurrentContext: viper.GetString("context"),
		},
	).ClientConfig()
	if err != nil {
		return true, fmt.Errorf("failed to create Kubernetes cluster configuration: %+v", err)
	}

	env := &envtest.Environment{
		Config: cfg,
		CRDDirectoryPaths: []string{
			path.Join(ModuleDirectory, "manifests/resources"),
		},
		AttachControlPlaneOutput: true,
		UseExistingCluster:       func(b bool) *bool { return &b }(true),
	}

	// End-to-end tests require an exclusive lock on the environment so, e.g.,
	// we don't try to simultaneously install Tekton in two different packages.
	release, err := Exclusive(ctx, LockEndToEndEnvironment)
	if err != nil {
		return true, fmt.Errorf("failed to acquire exclusive lock: %+v", err)
	}
	defer release()

	cfg, err = env.Start()
	if err != nil {
		return true, fmt.Errorf("failed to connect to Kubernetes cluster: %+v", err)
	}
	defer env.Stop()

	log.Println("connected to Kubernetes cluster", cfg.Host)

	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		return true, fmt.Errorf("failed to configure client resource discovery: %+v", err)
	}

	client, err := client.New(cfg, client.Options{
		Scheme: TestScheme,
		Mapper: mapper,
	})
	if err != nil {
		return true, fmt.Errorf("failed to configure client: %+v", err)
	}

	if viper.GetBool("label_nodes") {
		nodes := &corev1.NodeList{}
		if err := client.List(ctx, nodes); err != nil {
			return true, fmt.Errorf("failed to list nodes to label: %+v", err)
		}

		for _, node := range nodes.Items {
			for name, value := range obj.PipelineRunPodNodeSelector {
				node.GetLabels()[name] = value
			}

			if err := client.Update(ctx, &node); err != nil {
				return true, fmt.Errorf("failed to update node %s: %+v", node.GetName(), err)
			}
		}
	}

	ifc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return true, fmt.Errorf("failed to configure client: %+v", err)
	}

	tkc, err := tekton.NewForConfig(cfg)
	if err != nil {
		return true, fmt.Errorf("failed to configure Tekton client: %+v", err)
	}

	e := &EndToEndEnvironment{
		RESTConfig:              cfg,
		RESTMapper:              mapper,
		ControllerRuntimeClient: client,
		Interface:               ifc,
		TektonInterface:         tkc,
	}

	if viper.GetBool("install_environment") {
		for _, opt := range opts {
			if err := opt(ctx, e); err != nil {
				return true, err
			}
		}
	}

	fn(e)
	return true, nil
}

func WithEndToEndEnvironment(t *testing.T, fn func(e *EndToEndEnvironment), opts ...EndToEndEnvironmentOption) {
	enabled, err := doEndToEndEnvironment(fn, opts...)
	require.NoError(t, err)
	if !enabled {
		log.Println("end-to-end tests disabled by configuration")
	}
}

func RunEndToEnd(m *testing.M, fn func(e *EndToEndEnvironment), opts ...EndToEndEnvironmentOption) int {
	if enabled, err := doEndToEndEnvironment(fn, opts...); err != nil {
		log.Println(err)
		return 1
	} else if !enabled {
		log.Println("end-to-end tests disabled by configuration")
		return 0
	}

	return m.Run()
}
