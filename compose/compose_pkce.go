// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/pocket-id/fosite"
	"github.com/pocket-id/fosite/handler/oauth2"
	"github.com/pocket-id/fosite/handler/pkce"
)

// OAuth2PKCEFactory creates a PKCE handler.
func OAuth2PKCEFactory(config fosite.Configurator, storage fosite.Storage, strategy interface{}) interface{} {
	return &pkce.Handler{
		Strategy: strategy.(oauth2.AuthorizeCodeStrategyProvider),
		Storage:  storage.(pkce.PKCERequestStorageProvider),
		Config:   config,
	}
}
