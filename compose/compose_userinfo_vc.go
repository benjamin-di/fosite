// Copyright © 2025 Ory Corp
// Modified by contributors to the Pocket ID fork.
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/pocket-id/fosite"
	"github.com/pocket-id/fosite/handler/verifiable"
)

// OIDCUserinfoVerifiableCredentialFactory creates a verifiable credentials
// handler.
func OIDCUserinfoVerifiableCredentialFactory(config fosite.Configurator, storage fosite.Storage, strategy any) any {
	return &verifiable.Handler{
		NonceManagerProvider: storage.(verifiable.NonceManagerProvider),
		Config:               config,
	}
}
