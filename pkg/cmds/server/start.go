package server

import (
	"fmt"
	"io"
	"net"

	"github.com/etcd-manager/etcd-discovery/pkg/manager"
	"github.com/etcd-manager/etcd-discovery/pkg/server"
	genericoptions "github.com/etcd-manager/etcd-discovery/pkg/server/options"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

type DiscoveryServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
	StdOut             io.Writer
	StdErr             io.Writer
}

func NewDiscoveryServerOptions(out, errOut io.Writer) *DiscoveryServerOptions {
	o := &DiscoveryServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(),
		StdOut:             out,
		StdErr:             errOut,
	}
	o.RecommendedOptions.SecureServing.BindPort = 2381
	return o
}

func (o *DiscoveryServerOptions) AddFlags(fs *pflag.FlagSet) {
	o.RecommendedOptions.AddFlags(fs)
}

func (o DiscoveryServerOptions) Validate(args []string) error {
	var errors []error
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *DiscoveryServerOptions) Complete() error {
	return nil
}

func (o DiscoveryServerOptions) Config() (*server.Config, error) {
	// TODO have a "real" external address
	alternateIPs := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	if ip, err := o.RecommendedOptions.SecureServing.DefaultExternalAddress(); err == nil {
		alternateIPs = append(alternateIPs, ip)
	}
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, alternateIPs); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	config := &server.Config{
		GenericConfig: genericapiserver.NewRecommendedConfig(server.Codecs),
		EtcdConfig:    manager.NewEtcdConfig(),
	}
	if err := o.RecommendedOptions.ApplyTo(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (o DiscoveryServerOptions) Run(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	srv, err := config.Complete().New()
	if err != nil {
		return err
	}

	return srv.GenericAPIServer.PrepareRun().Run(stopCh)
}
