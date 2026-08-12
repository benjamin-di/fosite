// Copyright © 2026 Pocket ID
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/ory/fosite"
)

const authorizationResponseIssuer = "https://issuer.example.com/"

func TestWriteAuthorizeResponseIssuer(t *testing.T) {
	for _, responseMode := range []ResponseModeType{
		ResponseModeDefault,
		ResponseModeQuery,
		ResponseModeFragment,
		ResponseModeFormPost,
	} {
		t.Run(fmt.Sprintf("response_mode=%s", responseMode), func(t *testing.T) {
			oauth2 := &Fosite{Config: &Config{AuthorizationResponseIssuer: authorizationResponseIssuer}}
			request := newAuthorizationResponseIssuerRequest(t, responseMode)
			response := NewAuthorizeResponse()
			response.AddParameter("code", "authorization-code")
			response.AddParameter("iss", "https://untrusted.example.com/")
			recorder := httptest.NewRecorder()

			oauth2.WriteAuthorizeResponse(context.Background(), recorder, request, response)

			requireAuthorizationResponseIssuer(t, responseMode, recorder, authorizationResponseIssuer)
		})
	}
}

func TestWriteAuthorizeErrorIssuer(t *testing.T) {
	for _, responseMode := range []ResponseModeType{
		ResponseModeDefault,
		ResponseModeQuery,
		ResponseModeFragment,
		ResponseModeFormPost,
	} {
		t.Run(fmt.Sprintf("response_mode=%s", responseMode), func(t *testing.T) {
			oauth2 := &Fosite{Config: &Config{AuthorizationResponseIssuer: authorizationResponseIssuer}}
			request := newAuthorizationResponseIssuerRequest(t, responseMode)
			recorder := httptest.NewRecorder()

			oauth2.WriteAuthorizeError(context.Background(), recorder, request, ErrAccessDenied)

			requireAuthorizationResponseIssuer(t, responseMode, recorder, authorizationResponseIssuer)
		})
	}
}

func newAuthorizationResponseIssuerRequest(t *testing.T, responseMode ResponseModeType) *AuthorizeRequest {
	t.Helper()

	redirectURI, err := url.Parse("https://client.example.com/callback?existing=value")
	require.NoError(t, err)

	request := NewAuthorizeRequest()
	request.RedirectURI = redirectURI
	request.ResponseMode = responseMode
	request.State = "request-state"
	request.Client = &DefaultClient{RedirectURIs: []string{redirectURI.String()}}
	return request
}

func requireAuthorizationResponseIssuer(t *testing.T, responseMode ResponseModeType, recorder *httptest.ResponseRecorder, issuer string) {
	t.Helper()

	if responseMode == ResponseModeFormPost {
		body := recorder.Body.String()
		require.Contains(t, body, fmt.Sprintf(`<input type="hidden" name="iss" value="%s"/>`, issuer))
		require.Equal(t, 1, strings.Count(body, `name="iss"`))
		return
	}

	redirectURI, err := url.Parse(recorder.Header().Get("Location"))
	require.NoError(t, err)

	parameters := redirectURI.Query()
	if responseMode == ResponseModeFragment {
		parameters, err = url.ParseQuery(redirectURI.Fragment)
		require.NoError(t, err)
	}

	require.Equal(t, []string{issuer}, parameters["iss"])
}
