package provider

import (
	"fmt"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
)

type ProviderRegistry struct {
	Providers      map[string]ports.Provider
	OauthProviders map[string]ports.OAuthProvider
}

func (r *ProviderRegistry) GetProvider(name string) (ports.Provider, error) {
	p, ok := r.Providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidProvider, name)
	}
	return p, nil
}

func (r *ProviderRegistry) GetOauthProvider(name string) (ports.OAuthProvider, error) {
	p, ok := r.OauthProviders[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidProvider, name)
	}
	return p, nil
}
