package options

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/authentication/authenticatorfactory"
	"k8s.io/apiserver/pkg/server"
)

type ClientCertAuthenticationOptions struct {
	// ClientCA is the certificate bundle for all the signers that you'll recognize for incoming client certificates
	ClientCA string
}

func (s *ClientCertAuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ClientCA, "peer-trusted-ca-file", s.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the peer-trusted-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
}

// DelegatingAuthenticationOptions provides an easy way for composing API servers to delegate their authentication to
// the root kube API server.  The API federator will act as
// a front proxy and direction connections will be able to delegate to the core kube API server
type DelegatingAuthenticationOptions struct {
	// CacheTTL is the length of time that a token authentication answer will be cached.
	CacheTTL time.Duration

	ClientCert ClientCertAuthenticationOptions
}

func NewDelegatingAuthenticationOptions() *DelegatingAuthenticationOptions {
	return &DelegatingAuthenticationOptions{
		// very low for responsiveness, but high enough to handle storms
		CacheTTL:   10 * time.Second,
		ClientCert: ClientCertAuthenticationOptions{},
	}
}

func (s *DelegatingAuthenticationOptions) Validate() []error {
	return nil
}

func (s *DelegatingAuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&s.CacheTTL, "authentication-token-webhook-cache-ttl", s.CacheTTL,
		"The duration to cache responses from the webhook token authenticator.")

	s.ClientCert.AddFlags(fs)
}

func (s *DelegatingAuthenticationOptions) ApplyTo(c *server.Config) error {
	if s == nil {
		c.Authenticator = nil
		return nil
	}

	clientCA, err := s.getClientCA()
	if err != nil {
		return err
	}
	c, err = c.ApplyClientCert(clientCA.ClientCA)
	if err != nil {
		return fmt.Errorf("unable to load client CA file: %v", err)
	}

	cfg, err := s.ToAuthenticationConfig()
	if err != nil {
		return err
	}
	authenticator, securityDefinitions, err := cfg.New()
	if err != nil {
		return err
	}

	c.Authenticator = authenticator
	if c.OpenAPIConfig != nil {
		c.OpenAPIConfig.SecurityDefinitions = securityDefinitions
	}
	c.SupportsBasicAuth = false

	return nil
}

func (s *DelegatingAuthenticationOptions) ToAuthenticationConfig() (authenticatorfactory.DelegatingAuthenticatorConfig, error) {
	clientCA, err := s.getClientCA()
	if err != nil {
		return authenticatorfactory.DelegatingAuthenticatorConfig{}, err
	}

	ret := authenticatorfactory.DelegatingAuthenticatorConfig{
		Anonymous:    false,
		CacheTTL:     s.CacheTTL,
		ClientCAFile: clientCA.ClientCA,
	}
	return ret, nil
}

func (s *DelegatingAuthenticationOptions) getClientCA() (*ClientCertAuthenticationOptions, error) {
	if len(s.ClientCert.ClientCA) > 0 {
		return &s.ClientCert, nil
	}
	return nil, fmt.Errorf("cluster doesn't provide peer-trusted-ca-file")
}
