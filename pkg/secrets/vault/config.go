package vault

import "github.com/puppetlabs/horsehead/v2/logging"

// Config is used to configure vault clients with authentication,
// policies and paths for fetching secrets.
type Config struct {
	Addr string
	// The path to the service account token file
	K8sServiceAccountTokenPath string
	// The mount path to use when logging in using Kubernetes
	K8sAuthMountPath string
	// Optional token if not using kubernetes auth
	Token string
	// The role we should use when logging in
	Role string
	// The bucket path segment we are proxying requests for secrets for
	Bucket string
	// The engine to use to form paths from
	EngineMount           string
	ScopedSecretsPath     string
	ScopedConnectionsPath string
	Logger                logging.Logger
}
