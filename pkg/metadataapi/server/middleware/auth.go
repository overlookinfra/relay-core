package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/puppetlabs/leg/instrumentation/alerts/trackers"
	"github.com/puppetlabs/relay-core/pkg/authenticate"
	"github.com/puppetlabs/relay-core/pkg/manager/api"
	"github.com/puppetlabs/relay-core/pkg/manager/builder"
	"github.com/puppetlabs/relay-core/pkg/manager/configmap"
	"github.com/puppetlabs/relay-core/pkg/manager/memory"
	"github.com/puppetlabs/relay-core/pkg/manager/reject"
	"github.com/puppetlabs/relay-core/pkg/manager/service"
	"github.com/puppetlabs/relay-core/pkg/manager/vault"
	"github.com/puppetlabs/relay-core/pkg/model"
	"github.com/puppetlabs/relay-pls/pkg/plspb"
	"k8s.io/client-go/kubernetes"
)

// Credential represents a valid authentication request.
type Credential struct {
	Managers model.MetadataManagers
	Tags     []trackers.Tag
}

// Authenticator maps an HTTP request to a credential, if possible.
type Authenticator interface {
	// Authenticate performs the request mapping to a credential. If the request
	// cannot be verified but no other error occurs, this method returns nil.
	Authenticate(r *http.Request) (*Credential, error)
}

type KubernetesAuthenticatorClientFactoryFunc func(token string) (kubernetes.Interface, error)

type KubernetesAuthenticator struct {
	// The factory to produce a scoped token from a request.
	factory KubernetesAuthenticatorClientFactoryFunc

	// Client for using a Kubernetes pod-lookup intermediary instead of Bearer
	// request headers.
	kubernetesClient *authenticate.KubernetesInterface

	// Log Service
	logServiceClient plspb.LogClient

	// Uses Vault for token decryption (Kubernetes intermediary).
	vaultClient      *vaultapi.Client
	vaultTransitPath string
	vaultTransitKey  string

	// Uses Vault for JWT verification.
	vaultJWTAuthAddr string
	vaultJWTAuthPath string
	vaultJWTAuthRole string

	// Static keys to use for JWT verification.
	keys []interface{}
}

var _ Authenticator = &KubernetesAuthenticator{}

func (ka *KubernetesAuthenticator) intermediary(r *http.Request, mgrs *builder.MetadataBuilder) authenticate.Intermediary {
	if ka.kubernetesClient == nil {
		// In this case we expect the JWT to be specified in the Authorization
		// header as a Bearer.
		return authenticate.NewHTTPAuthorizationHeaderIntermediary(r)
	}

	// Extract IP from request to hand to Kubernetes for pod
	// authentication.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	// Go's HTTP server should always give us a valid RemoteAddr, and if it
	// fails the intermediary will bail on an empty IP anyway.
	ki := authenticate.NewKubernetesIntermediary(ka.kubernetesClient, net.ParseIP(host))

	if ka.vaultClient == nil {
		// Done. We assume the token will be set directly on the pod
		// annotation.
		return ki
	}

	// Otherwise we chain to the Vault client to let it decrypt the token.
	return ki.Chain(func(ctx context.Context, raw authenticate.Raw, md *authenticate.KubernetesIntermediaryMetadata) (authenticate.Intermediary, error) {
		mgrs.SetActionMetadata(memory.NewActionMetadataManager(&model.ActionMetadata{
			Image: md.Image,
		}))

		return authenticate.NewVaultTransitIntermediary(
			ka.vaultClient,
			ka.vaultTransitPath,
			ka.vaultTransitKey,
			string(raw),
			authenticate.VaultTransitIntermediaryWithContext(authenticate.VaultTransitNamespaceContext(string(md.NamespaceUID))),
		), nil
	})
}

func (ka *KubernetesAuthenticator) resolver(mgrs *builder.MetadataBuilder) authenticate.Resolver {
	var delegates []authenticate.Resolver

	if ka.vaultJWTAuthAddr != "" {
		delegates = append(delegates, authenticate.NewStubConfigVaultResolver(
			ka.vaultJWTAuthAddr,
			ka.vaultJWTAuthPath,
			authenticate.VaultResolverWithRole(ka.vaultJWTAuthRole),
			authenticate.VaultResolverWithInjector(authenticate.VaultResolverInjectorFunc(func(ctx context.Context, claims *authenticate.Claims, md *authenticate.VaultResolverMetadata) error {
				if claims.RelayVaultEnginePath == "" {
					return nil
				}

				base := vault.NewKVV2Client(md.VaultClient, claims.RelayVaultEnginePath)

				if claims.RelayVaultConnectionPath != "" {
					mgrs.SetConnections(vault.NewConnectionManager(base.In(claims.RelayVaultConnectionPath)))
				}

				if claims.RelayVaultSecretPath != "" {
					mgrs.SetSecrets(vault.NewSecretManager(base.In(claims.RelayVaultSecretPath)))
				}

				return nil
			})),
		))
	}

	for _, key := range ka.keys {
		delegates = append(delegates, authenticate.NewKeyResolver(key))
	}

	return authenticate.NewAnyResolver(delegates)
}

func (ka *KubernetesAuthenticator) injector(mgrs *builder.MetadataBuilder, tags *[]trackers.Tag) authenticate.Injector {
	return authenticate.InjectorFunc(func(ctx context.Context, claims *authenticate.Claims) error {
		client, err := ka.factory(claims.KubernetesServiceAccountToken)
		if err != nil {
			return err
		}

		immutableMap := configmap.NewClientConfigMap(client, claims.KubernetesNamespaceName, claims.RelayKubernetesImmutableConfigMapName)
		mutableMap := configmap.NewClientConfigMap(client, claims.KubernetesNamespaceName, claims.RelayKubernetesMutableConfigMapName)

		action := claims.Action()

		model.IfStep(action, func(step *model.Step) {
			// Only a step can work with parameters and outputs. Other actions
			// will get the default rejection manager.
			mgrs.SetParameters(configmap.NewParameterManager(immutableMap))
			mgrs.SetStepOutputs(configmap.NewStepOutputManager(step, mutableMap))
		})

		if claims.RelayEventAPIURL != nil {
			mgrs.SetEvents(api.NewEventManager(action, claims.RelayEventAPIURL.URL.String(), claims.RelayEventAPIToken))
		}

		mgrs.SetConditions(configmap.NewConditionManager(action, immutableMap))
		mgrs.SetEnvironment(configmap.NewEnvironmentManager(action, immutableMap))
		mgrs.SetSpec(configmap.NewSpecManager(action, immutableMap))
		mgrs.SetState(configmap.NewStateManager(action, mutableMap))

		logContext := ""
		switch action.Type().Singular {
		case model.ActionTypeTrigger.Singular:
			logContext = fmt.Sprintf("tenants/%s/triggers/%s", claims.RelayTenantID, claims.RelayName)
		case model.ActionTypeStep.Singular:
			logContext = fmt.Sprintf("tenants/%s/runs/%s/steps/%s", claims.RelayTenantID, claims.RelayRunID, claims.RelayName)
		}

		if ka.logServiceClient != nil && logContext != "" {
			mgrs.SetLogs(service.NewLogManager(ka.logServiceClient, logContext))
		} else {
			mgrs.SetLogs(reject.LogManager)
		}

		ts := []trackers.Tag{
			{Key: "relay.domain.id", Value: claims.RelayDomainID},
			{Key: "relay.tenant.id", Value: claims.RelayTenantID},
			{Key: "relay.run.id", Value: claims.RelayRunID},
			{Key: "relay.action.type", Value: action.Type().Singular},
			{Key: "relay.action.name", Value: claims.RelayName},
		}
		for _, tag := range ts {
			if tag.Value != "" {
				*tags = append(*tags, tag)
			}
		}

		return nil
	})
}

func (ka *KubernetesAuthenticator) Authenticate(r *http.Request) (*Credential, error) {
	mgrs := builder.NewMetadataBuilder()
	var tags []trackers.Tag

	auth := authenticate.NewAuthenticator(
		ka.intermediary(r, mgrs),
		ka.resolver(mgrs),
		authenticate.AuthenticatorWithInjector(ka.injector(mgrs, &tags)),
	)

	if ok, err := auth.Authenticate(r.Context()); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &Credential{
		Managers: mgrs.Build(),
		Tags:     tags,
	}, nil
}

type KubernetesAuthenticatorOption func(ka *KubernetesAuthenticator)

func KubernetesAuthenticatorWithKubernetesIntermediary(client *authenticate.KubernetesInterface) KubernetesAuthenticatorOption {
	return func(ka *KubernetesAuthenticator) {
		ka.kubernetesClient = client
	}
}

func KubernetesAuthenticatorWithLogServiceIntermediary(client plspb.LogClient) KubernetesAuthenticatorOption {
	return func(ka *KubernetesAuthenticator) {
		ka.logServiceClient = client
	}
}

func KubernetesAuthenticatorWithChainToVaultTransitIntermediary(client *vaultapi.Client, path, key string) KubernetesAuthenticatorOption {
	return func(ka *KubernetesAuthenticator) {
		ka.vaultClient = client
		ka.vaultTransitPath = path
		ka.vaultTransitKey = key
	}
}

func KubernetesAuthenticatorWithVaultResolver(addr, path, role string) KubernetesAuthenticatorOption {
	return func(ka *KubernetesAuthenticator) {
		ka.vaultJWTAuthAddr = addr
		ka.vaultJWTAuthPath = path
		ka.vaultJWTAuthRole = role
	}
}

func KubernetesAuthenticatorWithKeyResolver(key interface{}) KubernetesAuthenticatorOption {
	return func(ka *KubernetesAuthenticator) {
		ka.keys = append(ka.keys, key)
	}
}

func NewKubernetesAuthenticator(factory KubernetesAuthenticatorClientFactoryFunc, opts ...KubernetesAuthenticatorOption) *KubernetesAuthenticator {
	ka := &KubernetesAuthenticator{
		factory: factory,
	}

	for _, opt := range opts {
		opt(ka)
	}

	return ka
}

func WithAuthentication(a Authenticator) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cred, err := a.Authenticate(r)
			if err != nil {
				log(r.Context()).Error("failed to authenticate client", "error", err)
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
			} else if cred == nil {
				http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			} else {
				if capturer, ok := trackers.CapturerFromContext(r.Context()); ok {
					capturer = capturer.WithTags(cred.Tags...)
					r = r.WithContext(trackers.NewContextWithCapturer(r.Context(), capturer))
				}

				WithManagers(cred.Managers)(next).ServeHTTP(w, r)
			}
		})
	}
}
