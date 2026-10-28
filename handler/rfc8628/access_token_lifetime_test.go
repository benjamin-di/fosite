// Copyright © 2026 Pocket ID
// SPDX-License-Identifier: Apache-2.0

package rfc8628_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pocket-id/fosite"
	"github.com/pocket-id/fosite/handler/rfc8628"
	"github.com/pocket-id/fosite/storage"
)

func TestDeviceCodeAccessTokenLifetime(t *testing.T) {
	for _, phase := range []time.Duration{0, 499 * time.Millisecond, 500 * time.Millisecond, 999 * time.Millisecond} {
		t.Run(fmt.Sprint(phase), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				time.Sleep(phase)
				now := time.Now().UTC()
				ctx := context.Background()
				client := &fosite.DefaultClient{ID: "device-client", Public: true, GrantTypes: fosite.Arguments{string(fosite.GrantTypeDeviceCode)}}
				session := &fosite.DefaultSession{}
				session.SetExpiresAt(fosite.DeviceCode, now.Add(time.Minute))
				granted := &fosite.DeviceRequest{
					UserCodeState: fosite.UserCodeAccepted,
					Request: fosite.Request{
						ID: "device-grant", Client: client, Session: session, Form: url.Values{}, RequestedAt: now,
					},
				}
				code, signature, err := RFC8628HMACSHAStrategy.GenerateDeviceCode(ctx)
				require.NoError(t, err)
				store := storage.NewMemoryStore()
				require.NoError(t, store.CreateDeviceAuthSession(ctx, signature, "user-code-signature", granted))

				handler := rfc8628.DeviceCodeTokenEndpointHandler{
					Strategy: &mockDeviceCodeStrategyProvider{
						deviceRateLimitStrategy: &RFC8628HMACSHAStrategy,
						deviceCodeStrategy:      &RFC8628HMACSHAStrategy,
						userCodeStrategy:        &RFC8628HMACSHAStrategy,
						coreStrategy:            hmacshaStrategyOAuth,
					},
					Storage: store,
					Config: &fosite.Config{
						ScopeStrategy: fosite.HierarchicScopeStrategy, AudienceMatchingStrategy: fosite.DefaultAudienceMatchingStrategy,
						DeviceAndUserCodeLifespan: time.Minute, AccessTokenLifespan: 15 * time.Minute, RefreshTokenLifespan: time.Hour,
					},
				}
				request := fosite.NewAccessRequest(&fosite.DefaultSession{})
				request.GrantTypes = fosite.Arguments{string(fosite.GrantTypeDeviceCode)}
				request.Client = client
				request.Form = url.Values{"device_code": {code}}
				require.NoError(t, handler.HandleTokenEndpointRequest(ctx, request))

				expiresAt := request.GetSession().GetExpiresAt(fosite.AccessToken)
				require.Equal(t, int64(900), expiresAt.Unix()-now.Unix(), "expiry must not exceed the configured lifespan in Unix seconds")
				require.False(t, expiresAt.After(now.Add(15*time.Minute)))
			})
		})
	}
}
