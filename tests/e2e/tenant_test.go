package e2e_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	relayv1beta1 "github.com/puppetlabs/relay-core/pkg/apis/relay.sh/v1beta1"
	"github.com/puppetlabs/relay-core/pkg/model"
	"github.com/puppetlabs/relay-core/pkg/obj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestTenantFinalizer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	WithConfig(t, ctx, []ConfigOption{
		ConfigWithTenantReconciler,
	}, func(cfg *Config) {
		child := fmt.Sprintf("%s-child", cfg.Namespace.GetName())

		tenant := &relayv1beta1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "my-test-tenant",
			},
			Spec: relayv1beta1.TenantSpec{
				NamespaceTemplate: relayv1beta1.NamespaceTemplate{
					Metadata: metav1.ObjectMeta{
						Name: child,
					},
				},
			},
		}
		require.NoError(t, cfg.Environment.ControllerClient.Create(ctx, tenant))

		// Wait for namespace.
		require.NoError(t, retry.Wait(ctx, func(ctx context.Context) (bool, error) {
			if err := cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{
				Namespace: tenant.GetNamespace(),
				Name:      tenant.GetName(),
			}, tenant); err != nil {
				return true, err
			}

			for _, cond := range tenant.Status.Conditions {
				if cond.Type == relayv1beta1.TenantNamespaceReady && cond.Status == corev1.ConditionTrue {
					return true, nil
				}
			}

			return false, fmt.Errorf("waiting for namespace to be ready")
		}))

		// Get child namespace.
		namespace := &corev1.Namespace{}
		require.NoError(t, cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{Name: child}, namespace))

		// Delete tenant.
		require.NoError(t, cfg.Environment.ControllerClient.Delete(ctx, tenant))

		// Get child namespace again, should be gone after delete.
		require.NoError(t, retry.Wait(ctx, func(ctx context.Context) (bool, error) {
			if err := cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{Name: child}, namespace); errors.IsNotFound(err) {
				return true, nil
			} else if err != nil {
				return true, err
			}

			return false, fmt.Errorf("waiting for namespace to terminate")
		}))
	})
}

func TestTenantAPITriggerEventSinkMissingSecret(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	WithConfig(t, ctx, []ConfigOption{
		ConfigWithTenantReconciler,
	}, func(cfg *Config) {
		// Create tenant with event sink pointing at nonexistent secret.
		tenant := &relayv1beta1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "my-test-tenant",
			},
			Spec: relayv1beta1.TenantSpec{
				TriggerEventSink: relayv1beta1.TriggerEventSink{
					API: &relayv1beta1.APITriggerEventSink{
						URL: "http://stub.example.com",
						TokenFrom: &relayv1beta1.APITokenSource{
							SecretKeyRef: &relayv1beta1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "xyz",
								},
								Key: "test",
							},
						},
					},
				},
			},
		}
		require.NoError(t, cfg.Environment.ControllerClient.Create(ctx, tenant))

		// Wait for tenant to reconcile.
		var cond relayv1beta1.TenantCondition
		require.NoError(t, retry.Wait(ctx, func(ctx context.Context) (bool, error) {
			if err := cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{
				Namespace: tenant.GetNamespace(),
				Name:      tenant.GetName(),
			}, tenant); err != nil {
				return true, err
			}

			for _, cond = range tenant.Status.Conditions {
				if cond.Type == relayv1beta1.TenantEventSinkReady && cond.Status == corev1.ConditionFalse {
					return true, nil
				}
			}

			return false, fmt.Errorf("waiting for tenant to reconcile")
		}))
		assert.Equal(t, obj.TenantStatusReasonEventSinkNotConfigured, cond.Reason)
	})
}

func TestTenantAPITriggerEventSinkWithSecret(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	WithConfig(t, ctx, []ConfigOption{
		ConfigWithTenantReconciler,
	}, func(cfg *Config) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "my-test-tenant",
			},
			StringData: map[string]string{
				"token": "test",
			},
		}
		require.NoError(t, cfg.Environment.ControllerClient.Create(ctx, secret))

		tenant := &relayv1beta1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "my-test-tenant",
			},
			Spec: relayv1beta1.TenantSpec{
				TriggerEventSink: relayv1beta1.TriggerEventSink{
					API: &relayv1beta1.APITriggerEventSink{
						URL: "http://stub.example.com",
						TokenFrom: &relayv1beta1.APITokenSource{
							SecretKeyRef: &relayv1beta1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: secret.GetName(),
								},
								Key: "token",
							},
						},
					},
				},
			},
		}
		CreateAndWaitForTenant(t, ctx, cfg, tenant)
	})
}

func TestTenantAPITriggerEventSinkWithNamespaceAndSecret(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	WithConfig(t, ctx, []ConfigOption{
		ConfigWithTenantReconciler,
	}, func(cfg *Config) {
		child := fmt.Sprintf("%s-child", cfg.Namespace.GetName())

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "my-test-tenant",
			},
			StringData: map[string]string{
				"token": "test",
			},
		}
		require.NoError(t, cfg.Environment.ControllerClient.Create(ctx, secret))

		tenant := &relayv1beta1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "my-test-tenant",
			},
			Spec: relayv1beta1.TenantSpec{
				NamespaceTemplate: relayv1beta1.NamespaceTemplate{
					Metadata: metav1.ObjectMeta{
						Name: child,
					},
				},
				TriggerEventSink: relayv1beta1.TriggerEventSink{
					API: &relayv1beta1.APITriggerEventSink{
						URL: "http://stub.example.com",
						TokenFrom: &relayv1beta1.APITokenSource{
							SecretKeyRef: &relayv1beta1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: secret.GetName(),
								},
								Key: "token",
							},
						},
					},
				},
			},
		}
		CreateAndWaitForTenant(t, ctx, cfg, tenant)
	})
}

func TestTenantNamespaceUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	WithConfig(t, ctx, []ConfigOption{
		ConfigWithTenantReconciler,
	}, func(cfg *Config) {
		child1 := fmt.Sprintf("%s-child-1", cfg.Namespace.GetName())
		child2 := fmt.Sprintf("%s-child-2", cfg.Namespace.GetName())

		tenant := &relayv1beta1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "my-test-tenant",
			},
			Spec: relayv1beta1.TenantSpec{
				NamespaceTemplate: relayv1beta1.NamespaceTemplate{
					Metadata: metav1.ObjectMeta{
						Name: child1,
					},
				},
			},
		}
		CreateAndWaitForTenant(t, ctx, cfg, tenant)

		// Child namespace should now exist.
		var ns1 corev1.Namespace
		require.Equal(t, child1, tenant.Status.Namespace)
		require.NoError(t, cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{Name: child1}, &ns1))

		// Change namespace in tenant.
		Mutate(t, ctx, cfg, tenant, func() {
			tenant.Spec.NamespaceTemplate.Metadata.Name = child2
		})
		WaitForTenant(t, ctx, cfg, tenant)

		// First child namespace should now not exist or have deletion timestamp
		// set, second should exist.
		var ns2 corev1.Namespace
		require.Equal(t, child2, tenant.Status.Namespace)
		require.NoError(t, cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{Name: child2}, &ns2))

		if err := cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{Name: child1}, &ns1); err != nil {
			require.True(t, errors.IsNotFound(err))
		} else {
			require.NotEmpty(t, ns1.GetDeletionTimestamp())
		}
	})
}

func TestTenantToolInjection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	WithConfig(t, ctx, []ConfigOption{
		ConfigWithTenantReconciler,
	}, func(cfg *Config) {
		child := fmt.Sprintf("%s-child-1", cfg.Namespace.GetName())

		size, _ := resource.ParseQuantity("50Mi")
		storageClassName := "relay-hostpath"
		tenant := &relayv1beta1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cfg.Namespace.GetName(),
				Name:      "tenant-" + uuid.New().String(),
			},
			Spec: relayv1beta1.TenantSpec{
				NamespaceTemplate: relayv1beta1.NamespaceTemplate{
					Metadata: metav1.ObjectMeta{
						Name: child,
					},
				},
				ToolInjection: relayv1beta1.ToolInjection{
					VolumeClaimTemplate: &corev1.PersistentVolumeClaim{
						Spec: corev1.PersistentVolumeClaimSpec{
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceStorage: size,
								},
							},
							StorageClassName: &storageClassName,
						},
					},
				},
			},
		}

		CreateAndWaitForTenant(t, ctx, cfg, tenant)

		var ns corev1.Namespace
		require.Equal(t, child, tenant.Status.Namespace)
		require.NoError(t, cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{Name: child}, &ns))

		var pvc corev1.PersistentVolumeClaim
		require.NoError(t, cfg.Environment.ControllerClient.Get(ctx, client.ObjectKey{Name: tenant.GetName() + model.ToolInjectionVolumeClaimSuffixReadOnlyMany, Namespace: tenant.Status.Namespace}, &pvc))
		cfg.Environment.ControllerClient.Delete(ctx, &pvc)
	})
}
