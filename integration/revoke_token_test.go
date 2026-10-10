// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"net/http"
	"testing"

	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pocket-id/fosite"
	"github.com/pocket-id/fosite/compose"
	"github.com/pocket-id/fosite/handler/oauth2"
)

func TestRevokeToken(t *testing.T) {
	for _, strategy := range []oauth2.CoreStrategyProvider{
		hmacStrategyProvider,
	} {
		runRevokeTokenTest(t, strategy)
	}
}

func runRevokeTokenTest(t *testing.T, strategy oauth2.CoreStrategyProvider) {
	f := compose.Compose(new(fosite.Config), fositeStore, strategy, compose.OAuth2ClientCredentialsGrantFactory, compose.OAuth2TokenIntrospectionFactory, compose.OAuth2TokenRevocationFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2AppClient(ts)
	token, err := oauthClient.Token(t.Context())
	require.NoError(t, err)

	resp, _, errs := gorequest.New().Post(ts.URL+"/revoke").
		SetBasicAuth(oauthClient.ClientID, oauthClient.ClientSecret).
		Type("form").
		SendStruct(map[string]string{"token": "asdf"}).End()
	require.Len(t, errs, 0)
	assert.Equal(t, 200, resp.StatusCode)

	resp, _, errs = gorequest.New().Post(ts.URL+"/revoke").
		SetBasicAuth(oauthClient.ClientID, oauthClient.ClientSecret).
		Type("form").
		SendStruct(map[string]string{"token": token.AccessToken}).End()
	require.Len(t, errs, 0)
	assert.Equal(t, 200, resp.StatusCode)

	hres, _, errs := gorequest.New().Get(ts.URL+"/info").
		Set("Authorization", "bearer "+token.AccessToken).
		End()
	require.Len(t, errs, 0)
	assert.Equal(t, http.StatusUnauthorized, hres.StatusCode)
}
