// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"context"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/rfc8628"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/fosite/token/jwt"
)

type CommonStrategyProvider struct {
	CoreStrategy      oauth2.CoreStrategy
	AccessTokenStrat  oauth2.AccessTokenStrategy
	DeviceStrategy    *rfc8628.DefaultDeviceStrategy
	OIDCTokenStrategy openid.OpenIDConnectTokenStrategy
	jwt.Signer
}

var _ oauth2.AuthorizeCodeStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) AuthorizeCodeStrategy() oauth2.AuthorizeCodeStrategy {
	return s.CoreStrategy
}

var _ oauth2.AccessTokenStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) AccessTokenStrategy() oauth2.AccessTokenStrategy {
	if s.AccessTokenStrat != nil {
		return s.AccessTokenStrat
	}
	return s.CoreStrategy
}

var _ oauth2.RefreshTokenStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) RefreshTokenStrategy() oauth2.RefreshTokenStrategy {
	return s.CoreStrategy
}

var _ openid.OpenIDConnectTokenStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) OpenIDConnectTokenStrategy() openid.OpenIDConnectTokenStrategy {
	return s.OIDCTokenStrategy
}

var _ rfc8628.DeviceRateLimitStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) DeviceRateLimitStrategy() rfc8628.DeviceRateLimitStrategy {
	return s.DeviceStrategy
}

var _ rfc8628.DeviceCodeStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) DeviceCodeStrategy() rfc8628.DeviceCodeStrategy {
	return s.DeviceStrategy
}

var _ rfc8628.UserCodeStrategyProvider = (*CommonStrategyProvider)(nil)

func (s *CommonStrategyProvider) UserCodeStrategy() rfc8628.UserCodeStrategy {
	return s.DeviceStrategy
}

type HMACSHAStrategyConfigurator interface {
	fosite.AccessTokenLifespanProvider
	fosite.RefreshTokenLifespanProvider
	fosite.AuthorizeCodeLifespanProvider
	fosite.TokenEntropyProvider
	fosite.GlobalSecretProvider
	fosite.RotatedGlobalSecretsProvider
	fosite.HMACHashingProvider
	fosite.DeviceAndUserCodeLifespanProvider
}

func NewOAuth2HMACStrategy(config HMACSHAStrategyConfigurator) *oauth2.HMACSHAStrategy {
	return oauth2.NewHMACSHAStrategy(&hmac.HMACStrategy{Config: config}, config)
}

func NewOAuth2JWTStrategy(keyGetter func(context.Context) (interface{}, error), config fosite.Configurator) *oauth2.DefaultJWTStrategy {
	return &oauth2.DefaultJWTStrategy{
		Signer: &jwt.DefaultSigner{GetPrivateKey: keyGetter},
		Config: config,
	}
}

// NewOAuth2RFC9068JWTStrategy creates an access token strategy that issues RFC 9068 access tokens.
// Set it as CommonStrategyProvider.AccessTokenStrat so refresh tokens and authorization codes remain opaque.
func NewOAuth2RFC9068JWTStrategy(keyGetter func(context.Context) (interface{}, error), config fosite.Configurator) *oauth2.RFC9068JWTStrategy {
	return &oauth2.RFC9068JWTStrategy{
		DefaultJWTStrategy: NewOAuth2JWTStrategy(keyGetter, config),
	}
}

func NewOpenIDConnectStrategy(keyGetter func(context.Context) (interface{}, error), config fosite.Configurator) *openid.DefaultStrategy {
	return &openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{GetPrivateKey: keyGetter},
		Config: config,
	}
}

func NewDeviceStrategy(config fosite.Configurator) *rfc8628.DefaultDeviceStrategy {
	return &rfc8628.DefaultDeviceStrategy{
		Enigma: &hmac.HMACStrategy{Config: config},
		Config: config,
	}
}
