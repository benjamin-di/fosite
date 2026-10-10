// Copyright © 2026 Pocket ID
// SPDX-License-Identifier: Apache-2.0

package rfc8628_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/rfc8628"
	"github.com/ory/fosite/storage"
)

// usedDeviceCodeStore reports every device code as already used, like a
// storage backend that keeps invalidated device code sessions.
type usedDeviceCodeStore struct {
	*storage.MemoryStore
	used fosite.DeviceRequester
}

func (s *usedDeviceCodeStore) DeviceAuthStorage() rfc8628.DeviceAuthStorage {
	return s
}

func (s *usedDeviceCodeStore) GetDeviceCodeSession(context.Context, string, fosite.Session) (fosite.DeviceRequester, error) {
	return s.used, fosite.ErrInvalidatedDeviceCode
}

func TestDeviceCodeReuseRevokesTokensOfOriginalGrant(t *testing.T) {
	ctx := context.Background()
	client := &fosite.DefaultClient{ID: "foo", GrantTypes: fosite.Arguments{string(fosite.GrantTypeDeviceCode)}}

	// The first redemption of the device code issued tokens for the original grant
	original := &fosite.DeviceRequest{
		UserCodeState: fosite.UserCodeAccepted,
		Request: fosite.Request{
			ID:          "original-grant",
			Client:      client,
			Session:     &fosite.DefaultSession{},
			RequestedAt: time.Now().UTC(),
		},
	}
	memoryStore := storage.NewMemoryStore()
	require.NoError(t, memoryStore.CreateAccessTokenSession(ctx, "original-access-token", original))
	require.NoError(t, memoryStore.CreateRefreshTokenSession(ctx, "original-refresh-token", "original-access-token", original))

	h := rfc8628.DeviceCodeTokenEndpointHandler{
		Strategy: &mockDeviceCodeStrategyProvider{
			deviceRateLimitStrategy: &RFC8628HMACSHAStrategy,
			deviceCodeStrategy:      &RFC8628HMACSHAStrategy,
			userCodeStrategy:        &RFC8628HMACSHAStrategy,
			coreStrategy:            hmacshaStrategyOAuth,
		},
		Storage: &usedDeviceCodeStore{MemoryStore: memoryStore, used: original},
		Config: &fosite.Config{
			ScopeStrategy:             fosite.HierarchicScopeStrategy,
			AudienceMatchingStrategy:  fosite.DefaultAudienceMatchingStrategy,
			DeviceAndUserCodeLifespan: time.Minute,
		},
	}

	// The replayed device code arrives with a new request ID
	code, _, err := RFC8628HMACSHAStrategy.GenerateDeviceCode(ctx)
	require.NoError(t, err)
	replay := fosite.NewAccessRequest(&fosite.DefaultSession{})
	replay.GrantTypes = fosite.Arguments{string(fosite.GrantTypeDeviceCode)}
	replay.Client = client
	replay.Form = url.Values{"device_code": {code}}
	require.NotEqual(t, original.GetID(), replay.GetID())

	err = h.HandleTokenEndpointRequest(ctx, replay)
	require.ErrorIs(t, err, fosite.ErrInvalidGrant)

	// Both tokens of the original grant are revoked
	_, err = memoryStore.GetAccessTokenSession(ctx, "original-access-token", nil)
	require.ErrorIs(t, err, fosite.ErrNotFound)
	_, err = memoryStore.GetRefreshTokenSession(ctx, "original-refresh-token", nil)
	require.ErrorIs(t, err, fosite.ErrInactiveToken)
}
