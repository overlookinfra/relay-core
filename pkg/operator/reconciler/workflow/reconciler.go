package workflow

import (
	"context"
	"fmt"
	"io"

	"github.com/puppetlabs/horsehead/v2/storage"
	"github.com/puppetlabs/relay-core/pkg/authenticate"
	"github.com/puppetlabs/relay-core/pkg/errmark"
	"github.com/puppetlabs/relay-core/pkg/operator/dependency"
	"github.com/puppetlabs/relay-core/pkg/operator/obj"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
	"knative.dev/pkg/apis"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	*dependency.DependencyManager

	Client client.Client
	Scheme *runtime.Scheme

	standalone bool
	metrics    *controllerObservations
	issuer     authenticate.Issuer
}

func NewReconciler(dm *dependency.DependencyManager) *Reconciler {
	return &Reconciler{
		DependencyManager: dm,

		Client: dm.Manager.GetClient(),
		Scheme: dm.Manager.GetScheme(),

		standalone: dm.Config.Standalone,
		metrics:    newControllerObservations(dm.Metrics),
		issuer: authenticate.IssuerFunc(func(ctx context.Context, claims *authenticate.Claims) (authenticate.Raw, error) {
			raw, err := authenticate.NewKeySignerIssuer(dm.JWTSigner).Issue(ctx, claims)
			if err != nil {
				return nil, err
			}

			return authenticate.NewVaultTransitWrapper(
				dm.VaultClient,
				dm.Config.VaultTransitPath,
				dm.Config.VaultTransitKey,
				authenticate.VaultTransitWrapperWithContext(authenticate.VaultTransitNamespaceContext(claims.KubernetesNamespaceUID)),
			).Wrap(ctx, raw)
		}),
	}
}

func (r *Reconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	ctx := context.Background()

	wr := obj.NewWorkflowRun(req.NamespacedName)
	if ok, err := wr.Load(ctx, r.Client); err != nil {
		return ctrl.Result{}, errmark.MapLast(err, func(err error) error {
			return fmt.Errorf("failed to load dependencies: %+v", err)
		})
	} else if !ok {
		// CRD deleted from under us?
		return ctrl.Result{}, nil
	}

	if ts := wr.Object.GetDeletionTimestamp(); ts != nil && !ts.IsZero() {
		return ctrl.Result{}, nil
	}

	if len(wr.Object.Spec.Workflow.Steps) == 0 {
		if err := wr.Complete(ctx, r.Client); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	var pr *obj.PipelineRun
	err = r.metrics.trackDurationWithOutcome(metricWorkflowRunStartUpDuration, func() error {
		deps, err := obj.ApplyWorkflowRunDeps(
			ctx,
			r.Client,
			wr,
			r.issuer,
			r.Config.MetadataAPIURL,
			obj.WorkflowRunDepsWithStandaloneMode(r.standalone),
		)

		if err != nil {
			err = errmark.MarkTransient(err, obj.TransientIfRequired)

			return errmark.MapLast(err, func(err error) error {
				return fmt.Errorf("failed to apply dependencies: %+v", err)
			})
		}

		pipeline, err := obj.ApplyPipeline(ctx, r.Client, deps)
		if err != nil {
			return errmark.MapLast(err, func(err error) error {
				return fmt.Errorf("failed to apply Pipeline: %+v", err)
			})
		}

		pr, err = obj.ApplyPipelineRun(ctx, r.Client, pipeline)
		if err != nil {
			return errmark.MapLast(err, func(err error) error {
				return fmt.Errorf("failed to apply PipelineRun: %+v", err)
			})
		}

		return nil
	})
	if err != nil {
		var errUser bool
		errmark.IfUser(errmark.Resolve(err), func(err error) {
			errUser = true
		})

		if errUser {
			klog.Error(err)
			if err := wr.Fail(ctx, r.Client); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	err = r.metrics.trackDurationWithOutcome(metricWorkflowRunLogUploadDuration, func() error {
		r.uploadLogs(ctx, wr, pr)
		return nil
	})
	if err != nil {
		klog.Warning(err)
	}

	obj.ConfigureWorkflowRun(wr, pr)

	if err := wr.PersistStatus(ctx, r.Client); err != nil {
		return ctrl.Result{}, errmark.MapLast(err, func(err error) error {
			return fmt.Errorf("failed to persist WorkflowRun: %+v", err)
		})
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) uploadLogs(ctx context.Context, wr *obj.WorkflowRun, plr *obj.PipelineRun) {
	podNames := make(map[string]string)

	for name, tr := range plr.Object.Status.TaskRuns {
		if tr.Status == nil {
			continue
		}

		// TODO: This should support retries, possibly to different log
		// endpoints?
		if cond := tr.Status.GetCondition(apis.ConditionSucceeded); cond == nil || cond.IsUnknown() {
			continue
		}

		podNames[name] = tr.Status.PodName
	}

	for name, step := range wr.Object.Status.Steps {
		if step.LogKey != "" {
			// Already uploaded.
			continue
		}

		podName, found := podNames[step.Name]
		if !found {
			// Not done yet.
			klog.Infof("WorkflowRun %s step %q is still progressing, waiting to upload logs", wr.Key, name)
			continue
		}

		klog.Infof("WorkflowRun %s step %q is complete, uploading logs for pod %s", wr.Key, name, podName)

		logKey, err := r.uploadLog(ctx, plr.Key.Namespace, podName, "step-step")
		if err != nil {
			klog.Warningf("failed to upload log for WorkflowRun %s step %q: %+v", wr.Key, name, err)
		}

		step.LogKey = logKey
		wr.Object.Status.Steps[name] = step
	}
}

func (r *Reconciler) uploadLog(ctx context.Context, namespace string, podName string, containerName string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", namespace, podName, containerName)

	// XXX: We can't do this with the dynamic client yet.
	opts := &corev1.PodLogOptions{
		Container: containerName,
	}
	rc, err := r.KubeClient.CoreV1().Pods(namespace).GetLogs(podName, opts).Stream()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	storageOpts := storage.PutOptions{
		ContentType: "application/octet-stream",
	}

	err = r.StorageClient.Put(ctx, key, func(w io.Writer) error {
		_, err := io.Copy(w, rc)

		return err
	}, storageOpts)
	if err != nil {
		return "", err
	}

	return key, nil
}
