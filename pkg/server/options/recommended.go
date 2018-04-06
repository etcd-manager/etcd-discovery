package options

import (
	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

// RecommendedOptions contains the recommended options for running an API server.
// If you add something to this list, it should be in a logical grouping.
// Each of them can be nil to leave the feature unconfigured on ApplyTo.
type RecommendedOptions struct {
	SecureServing  *genericoptions.SecureServingOptions
	Authentication *DelegatingAuthenticationOptions
	Audit          *genericoptions.AuditOptions
	Features       *genericoptions.FeatureOptions
}

func NewRecommendedOptions() *RecommendedOptions {
	return &RecommendedOptions{
		SecureServing:  genericoptions.NewSecureServingOptions(),
		Authentication: NewDelegatingAuthenticationOptions(),
		Audit:          genericoptions.NewAuditOptions(),
		Features:       genericoptions.NewFeatureOptions(),
	}
}

func (o *RecommendedOptions) AddFlags(fs *pflag.FlagSet) {
	o.SecureServing.AddFlags(fs)
	o.Authentication.AddFlags(fs)
	o.Audit.AddFlags(fs)
	o.Features.AddFlags(fs)
}

func (o *RecommendedOptions) ApplyTo(config *server.RecommendedConfig) error {
	if err := o.SecureServing.ApplyTo(&config.Config); err != nil {
		return err
	}
	if err := o.Authentication.ApplyTo(&config.Config); err != nil {
		return err
	}
	if err := o.Audit.ApplyTo(&config.Config); err != nil {
		return err
	}
	if err := o.Features.ApplyTo(&config.Config); err != nil {
		return err
	}
	return nil
}

func (o *RecommendedOptions) Validate() []error {
	var errors []error
	errors = append(errors, o.SecureServing.Validate()...)
	errors = append(errors, o.Authentication.Validate()...)
	errors = append(errors, o.Audit.Validate()...)
	errors = append(errors, o.Features.Validate()...)

	return errors
}
