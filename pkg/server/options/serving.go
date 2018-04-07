package options

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"

	"github.com/appscode/kutil/tools/certstore"
	"github.com/golang/glog"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/authentication/authenticatorfactory"
	"k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/util/cert"
	certutil "k8s.io/client-go/util/cert"
)

const organization = "system:etcd"

type SecureServingOptions struct {
	BindAddress net.IP
	BindPort    int
	// BindNetwork is the type of network to bind to - defaults to "tcp", accepts "tcp",
	// "tcp4", and "tcp6".
	BindNetwork string

	// Listener is the secure server network listener.
	// either Listener or BindAddress/BindPort/BindNetwork is set,
	// if Listener is set, use it and omit BindAddress/BindPort/BindNetwork.
	Listener net.Listener

	// CertDirectory is a directory that will contain the certificates.  If the cert and key aren't specifically set
	// this will be used to derive a match with the "pair-name"
	CertDirectory string

	// PeerCert is the TLS cert info for serving secure peer (server-to-server / cluster) traffic
	PeerCert GeneratableKeyCert

	// ServerCert is the TLS cert info for serving secure client-to-server traffic
	ServerCert GeneratableKeyCert
}

type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string
}

type GeneratableKeyCert struct {
	CertKey CertKey

	// CACertFile is an optional file containing the certificate chain for CertKey.CertFile
	CACertFile string
	// PairName is the name which will be used with CertDirectory to make a cert and key names
	// It becomes CertDirector/PairName.crt and CertDirector/PairName.key
	PairName string
}

func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress:   net.ParseIP("0.0.0.0"),
		BindPort:      443,
		CertDirectory: "etcd.local.config/certificates",
		PeerCert: GeneratableKeyCert{
			PairName: "peer",
		},
		ServerCert: GeneratableKeyCert{
			PairName: "server",
		},
	}
}

func (s *SecureServingOptions) DefaultExternalAddress() (net.IP, error) {
	return utilnet.ChooseBindAddress(s.BindAddress)
}

func (s *SecureServingOptions) Validate() []error {
	if s == nil {
		return nil
	}

	errors := []error{}

	if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(errors, fmt.Errorf("--secure-port %v must be between 0 and 65535, inclusive. 0 for turning off secure port.", s.BindPort))
	}

	return errors
}

func (s *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	if s == nil {
		return
	}

	fs.IPVar(&s.BindAddress, "bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --secure-port port. The "+
		"associated interface(s) must be reachable by the rest of the cluster, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0).")

	fs.IntVar(&s.BindPort, "secure-port", s.BindPort, ""+
		"The port on which to serve HTTPS with authentication and authorization. If 0, "+
		"don't serve HTTPS at all.")

	fs.StringVar(&s.CertDirectory, "cert-dir", s.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --peer-cert-file and --peer-private-key-file are provided, this flag will be ignored.")

	fs.StringVar(&s.PeerCert.CertKey.CertFile, "peer-cert-file", s.PeerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate used for SSL/TLS connections between peers. "+
		"This will be used both for listening on the peer address as well as sending requests to "+
		"other peers. If HTTPS serving is enabled, and --peer-cert-file and --peer-private-key-file "+
		"are not provided, a self-signed certificate and key are generated for the public address and "+
		"saved to the directory specified by --cert-dir.")

	fs.StringVar(&s.PeerCert.CertKey.KeyFile, "peer-private-key-file", s.PeerCert.CertKey.KeyFile,
		"File containing the default x509 private key matching --peer-cert-file.")

	fs.StringVar(&s.PeerCert.CACertFile, "peer-trusted-ca-file", s.PeerCert.CACertFile, ""+
		"File containing the certificate authority will used for secure access from peer etcd servers. "+
		"This must be a valid PEM-encoded CA bundle.")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate used for SSL/TLS connections to etcd. When this "+
		"option is set, advertise-client-urls can use the HTTPS schema. If HTTPS serving is enabled, "+
		"and --cert-file and --private-key-file are not provided, a self-signed certificate and key are "+
		"generated for the public address and saved to the directory specified by --cert-dir.")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "private-key-file", s.ServerCert.CertKey.KeyFile,
		"File containing the default x509 private key matching --cert-file.")

	fs.StringVar(&s.ServerCert.CACertFile, "trusted-ca-file", s.ServerCert.CACertFile, ""+
		"File containing the certificate authority will used for secure client-to-server communication. "+
		"This must be a valid PEM-encoded CA bundle.")
}

// ApplyTo fills up serving information in the server configuration.
func (s *SecureServingOptions) ApplyTo(c *server.Config) error {
	if s == nil {
		c.Authenticator = nil
		return nil
	}
	if s.BindPort <= 0 {
		return nil
	}

	if s.Listener == nil {
		var err error
		addr := net.JoinHostPort(s.BindAddress.String(), strconv.Itoa(s.BindPort))
		s.Listener, s.BindPort, err = genericoptions.CreateListener(s.BindNetwork, addr)
		if err != nil {
			return fmt.Errorf("failed to create listener: %v", err)
		}
	}

	if err := s.applyServingInfoTo(c); err != nil {
		return err
	}

	c.SecureServingInfo.Listener = s.Listener

	// create self-signed cert+key with the fake server.LoopbackClientServerNameOverride and
	// let the server return it when the loopback client connects.
	certPem, keyPem, err := certutil.GenerateSelfSignedCertKey(server.LoopbackClientServerNameOverride, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to generate self-signed certificate for loopback connection: %v", err)
	}
	tlsCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return fmt.Errorf("failed to generate self-signed certificate for loopback connection: %v", err)
	}

	secureLoopbackClientConfig, err := c.SecureServingInfo.NewLoopbackClientConfig(uuid.NewRandom().String(), certPem)
	switch {
	// if we failed and there's no fallback loopback client config, we need to fail
	case err != nil && c.LoopbackClientConfig == nil:
		return err

		// if we failed, but we already have a fallback loopback client config (usually insecure), allow it
	case err != nil && c.LoopbackClientConfig != nil:

	default:
		c.LoopbackClientConfig = secureLoopbackClientConfig
		c.SecureServingInfo.SNICerts[server.LoopbackClientServerNameOverride] = &tlsCert
	}

	return nil
}

func (s *SecureServingOptions) applyServingInfoTo(c *server.Config) error {
	if len(s.PeerCert.CACertFile) == 0 {
		return fmt.Errorf("cluster doesn't provide --peer-ca-file")
	}

	secureServingInfo := &server.SecureServingInfo{}

	serverCertFile, serverKeyFile := s.PeerCert.CertKey.CertFile, s.PeerCert.CertKey.KeyFile
	// load main cert
	if len(serverCertFile) != 0 || len(serverKeyFile) != 0 {
		tlsCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
		if err != nil {
			return fmt.Errorf("unable to load server certificate: %v", err)
		}
		secureServingInfo.Cert = &tlsCert
	}

	// load CA cert
	pemData, err := ioutil.ReadFile(s.PeerCert.CACertFile)
	if err != nil {
		return fmt.Errorf("failed to read certificate authority from %q: %v", s.PeerCert.CACertFile, err)
	}
	block, pemData := pem.Decode(pemData)
	if block == nil {
		return fmt.Errorf("no certificate found in certificate authority file %q", s.PeerCert.CACertFile)
	}
	if block.Type != "CERTIFICATE" {
		return fmt.Errorf("expected CERTIFICATE block in certiticate authority file %q, found: %s", s.PeerCert.CACertFile, block.Type)
	}
	secureServingInfo.CACert = &tls.Certificate{
		Certificate: [][]byte{block.Bytes},
	}
	secureServingInfo.SNICerts = map[string]*tls.Certificate{}

	c.SecureServingInfo = secureServingInfo
	c.ReadWritePort = s.BindPort

	// require client cert auth
	c, err = c.ApplyClientCert(s.PeerCert.CACertFile)
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

func (s *SecureServingOptions) MaybeDefaultWithSelfSignedCerts(publicAddress string, alternateDNS []string, alternateIPs []net.IP) error {
	if s == nil {
		return nil
	}
	keyCert := &s.PeerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	// add either the bind address or localhost to the valid alternates
	bindIP := s.BindAddress.String()
	if bindIP == "0.0.0.0" {
		alternateDNS = append(alternateDNS, "localhost")
	} else {
		alternateIPs = append(alternateIPs, s.BindAddress)
	}
	sans := cert.AltNames{
		IPs:      alternateIPs,
		DNSNames: alternateDNS,
	}

	// peer certs
	err := s.generatePeerCerts(publicAddress, sans)
	if err != nil {
		return fmt.Errorf("unable to generate self signed peer cert: %v", err)
	} else {
		glog.Infof("Generated self-signed peer cert (%s, %s)", keyCert.CertFile, keyCert.KeyFile)
	}

	// server certs
	err = s.generateServerCerts(publicAddress, sans)
	if err != nil {
		return fmt.Errorf("unable to generate self signed server cert: %v", err)
	} else {
		glog.Infof("Generated self-signed server cert (%s, %s)", keyCert.CertFile, keyCert.KeyFile)
	}

	return nil
}

func (s *SecureServingOptions) generatePeerCerts(host string, sans cert.AltNames) error {
	store, err := certstore.NewCertStore(afero.NewOsFs(), s.CertDirectory, organization)
	if err != nil {
		return errors.Wrap(err, "failed to create certificate store.")
	}
	err = store.InitCA(s.PeerCert.PairName)
	if err != nil {
		return errors.Wrap(err, "failed to init ca.")
	}
	crt, key, err := store.NewPeerCertPair(host, sans)
	if err != nil {
		return err
	}
	err = store.WriteBytes(host, crt, key)
	if err != nil {
		return err
	}

	s.PeerCert.CACertFile = store.CertFile(store.CAName())
	s.PeerCert.CertKey.CertFile = store.CertFile(host)
	s.PeerCert.CertKey.KeyFile = store.KeyFile(host)
	return nil
}

func (s *SecureServingOptions) generateServerCerts(host string, sans cert.AltNames) error {
	store, err := certstore.NewCertStore(afero.NewOsFs(), s.CertDirectory, organization)
	if err != nil {
		return errors.Wrap(err, "failed to create certificate store.")
	}
	err = store.InitCA(s.ServerCert.PairName)
	if err != nil {
		return errors.Wrap(err, "failed to init ca.")
	}
	crt, key, err := store.NewServerCertPair(host, sans)
	if err != nil {
		return err
	}
	err = store.WriteBytes(host, crt, key)
	if err != nil {
		return err
	}

	s.ServerCert.CACertFile = store.CertFile(store.CAName())
	s.ServerCert.CertKey.CertFile = store.CertFile(host)
	s.ServerCert.CertKey.KeyFile = store.KeyFile(host)
	return nil
}

func (s *SecureServingOptions) ToAuthenticationConfig() (authenticatorfactory.DelegatingAuthenticatorConfig, error) {
	ret := authenticatorfactory.DelegatingAuthenticatorConfig{
		Anonymous:    false,
		CacheTTL:     0,
		ClientCAFile: s.PeerCert.CACertFile,
	}
	return ret, nil
}
