// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/errorsx"
	"github.com/pocket-id/fosite"
	"github.com/pocket-id/fosite/token/jwt"
)

// DefaultJWTStrategy is a JWT RS256 strategy.
type DefaultJWTStrategy struct {
	Signer jwt.Signer
	Config interface {
		fosite.AccessTokenIssuerProvider
		fosite.JWTScopeFieldProvider
	}
}

var _ AccessTokenStrategy = (*DefaultJWTStrategy)(nil)

func signature(token string) string {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return ""
	}

	return split[2]
}

func (h *DefaultJWTStrategy) AccessTokenSignature(_ context.Context, token string) string {
	return signature(token)
}

func (h *DefaultJWTStrategy) GenerateAccessToken(ctx context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(ctx, fosite.AccessToken, requester)
}

func (h *DefaultJWTStrategy) ValidateAccessToken(ctx context.Context, _ fosite.Requester, token string) error {
	_, err := validate(ctx, h.Signer, token)
	return err
}

func validate(ctx context.Context, jwtStrategy jwt.Signer, token string) (t *jwt.Token, err error) {
	t, err = jwtStrategy.Decode(ctx, token)
	if err == nil {
		err = t.Claims.Valid()
		return
	}

	var e *jwt.ValidationError
	if errors.As(err, &e) {
		err = errorsx.WithStack(toRFCErr(e).WithWrap(err).WithDebug(err.Error()))
	}

	return
}

func toRFCErr(v *jwt.ValidationError) *fosite.RFC6749Error {
	switch {
	case v == nil:
		return nil
	case v.Has(jwt.ValidationErrorMalformed):
		return fosite.ErrInvalidTokenFormat
	case v.Has(jwt.ValidationErrorUnverifiable | jwt.ValidationErrorSignatureInvalid):
		return fosite.ErrTokenSignatureMismatch
	case v.Has(jwt.ValidationErrorExpired):
		return fosite.ErrTokenExpired
	case v.Has(jwt.ValidationErrorAudience |
		jwt.ValidationErrorIssuedAt |
		jwt.ValidationErrorIssuer |
		jwt.ValidationErrorNotValidYet |
		jwt.ValidationErrorId |
		jwt.ValidationErrorClaimsInvalid):
		return fosite.ErrTokenClaim
	default:
		return fosite.ErrRequestUnauthorized
	}
}

func (h *DefaultJWTStrategy) generate(ctx context.Context, tokenType fosite.TokenType, requester fosite.Requester) (string, string, error) {
	claims, header, err := h.prepareJWT(ctx, tokenType, requester)
	if err != nil {
		return "", "", err
	}

	return h.Signer.Generate(ctx, claims, header)
}

func (h *DefaultJWTStrategy) prepareJWT(ctx context.Context, tokenType fosite.TokenType, requester fosite.Requester) (jwt.MapClaims, *jwt.Headers, error) {
	if jwtSession, ok := requester.GetSession().(JWTSessionContainer); !ok {
		return nil, nil, errors.Errorf("Session must be of type JWTSessionContainer but got type: %T", requester.GetSession())
	} else if claims := jwtSession.GetJWTClaims(); claims == nil {
		return nil, nil, errors.New("GetTokenClaims() must not be nil")
	} else {
		mapClaims := claims.
			With(
				time.Now().UTC(),
				jwtSession.GetExpiresAt(tokenType),
				requester.GetGrantedScopes(),
				requester.GetGrantedAudience(),
			).
			WithDefaults(
				h.Config.GetAccessTokenIssuer(ctx),
			).
			WithScopeField(
				h.Config.GetJWTScopeField(ctx),
			).
			ToMapClaims()

		return mapClaims, jwtSession.GetJWTHeader(), nil
	}
}
