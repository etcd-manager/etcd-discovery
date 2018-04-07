package options

import (
	"github.com/etcd-manager/etcd-discovery/pkg/server"
	"github.com/spf13/pflag"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

// RecommendedOptions contains the recommended options for running an API server.
// If you add something to this list, it should be in a logical grouping.
// Each of them can be nil to leave the feature unconfigured on ApplyTo.
type RecommendedOptions struct {
	Etcd          *EtcdOptions
	SecureServing *SecureServingOptions
	Audit         *genericoptions.AuditOptions
	Features      *genericoptions.FeatureOptions
}

func NewRecommendedOptions() *RecommendedOptions {
	return &RecommendedOptions{
		Etcd:          NewEtcdOptions(),
		SecureServing: NewSecureServingOptions(),
		Audit:         genericoptions.NewAuditOptions(),
		Features:      genericoptions.NewFeatureOptions(),
	}
}

func (o *RecommendedOptions) AddFlags(fs *pflag.FlagSet) {
	o.Etcd.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	o.Audit.AddFlags(fs)
	o.Features.AddFlags(fs)
}

func (o *RecommendedOptions) ApplyTo(config *server.Config) error {
	if err := o.Etcd.ApplyTo(config.ExtraConfig); err != nil {
		return err
	}
	if err := o.SecureServing.ApplyTo(&config.GenericConfig.Config); err != nil {
		return err
	}
	if err := o.Audit.ApplyTo(&config.GenericConfig.Config); err != nil {
		return err
	}
	if err := o.Features.ApplyTo(&config.GenericConfig.Config); err != nil {
		return err
	}
	return nil
}

func (o *RecommendedOptions) Validate() []error {
	var errors []error
	errors = append(errors, o.SecureServing.Validate()...)
	errors = append(errors, o.Audit.Validate()...)
	errors = append(errors, o.Features.Validate()...)
	return errors
}
